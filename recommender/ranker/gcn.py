import os

import numpy as np
import pandas as pd
import torch
import torch_geometric.transforms as T
from data import TrainData
from ranker.ranker import Ranker
from torch_geometric.data import HeteroData
from torch_geometric.nn import HeteroConv, SAGEConv
from user import User

RND_STATE = 42
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


def bpr_loss(pred, target):
    """
    pred   : Tensor [N]
    target : Tensor [N] with values {0, 1}
    """

    pred = pred.view(-1)
    target = target.view(-1)

    pos_pred = pred[target == 1]
    neg_pred = pred[target == 0]

    # no valid pairs
    if pos_pred.numel() == 0 or neg_pred.numel() == 0:
        return torch.tensor(0.0, device=pred.device, requires_grad=True)

    # sample one negative per positive
    neg_idx = torch.randint(
        0,
        neg_pred.size(0),
        (pos_pred.size(0),),
        device=pred.device,
    )

    neg_sampled = neg_pred[neg_idx]

    loss = -torch.log(torch.sigmoid(pos_pred - neg_sampled)).mean()

    return loss


class GCNRanker(Ranker):
    def __init__(self, train_data: TrainData):
        super().__init__(train_data)
        self.set_train_params()

    def fit(self) -> None:
        self.id_to_povinn = dict(
            zip(self.train_data.povinn["course_id"], self.train_data.povinn["povinn"])
        )

        neg_train, neg_val, _ = negative_split(
            self.train_data.train,
            self.train_data.val,
            None,
            self.train_data.finished,
            val_ratio=-1,
            train_ratio=1,
        )
        self.train_graph, self.val_graph, _ = graph_data(
            self.train_data.user,
            self.train_data.povinn,
            self.train_data.train,
            self.train_data.val,
            None,
            neg_train,
            neg_val,
            None,
        )

        self.model = Model(self.train_graph, hidden_channels=64, out_channels=32).to(
            device
        )

        if os.path.exists("gcn_ranker.pth"):
            self.model.load_state_dict(torch.load("gcn_ranker.pth", weights_only=True))
        else:
            optimizer = torch.optim.AdamW(self.model.parameters(), lr=self.lr)
            train_model(
                self.model,
                optimizer,
                self.train_graph,
                self.val_graph,
                self.loss_fn,
                self.epochs,
            )

            u = User(self.train_data.rand_soident_from_dev(), "", 0, None)
            u.fetch = True
            _ = self.rank(u)  # init lazy loaded params

            torch.save(self.model.state_dict(), "gcn_ranker.pth")

        self.model.eval()

    def rank(self, user: User) -> list[str]:
        if user.fetch:
            uid = self.train_data.user[self.train_data.user["soident"] == int(user.id)][
                "user_id"
            ].iloc[0]
            graph = self.val_graph
        else:
            uid, graph = self.graph_with_user(user)

        with torch.no_grad():
            pred = self.model(
                graph.x_dict,
                graph.edge_index_dict,
                graph["user", "course"].edge_label_index,
            )
        user_id = graph["user", "course"].edge_label_index[0].cpu().numpy()
        course_id = graph["user", "course"].edge_label_index[1].cpu().numpy()
        results = pd.DataFrame(
            {
                "user_id": user_id,
                "course_id": course_id,
                "pred": pred.sigmoid().cpu().numpy(),
            }
        )
        results = results[results["user_id"] == uid]
        results = results.groupby("user_id").agg({"pred": list, "course_id": list})
        if results.empty:
            pred = []
        else:
            pred_ranking = results["pred"].iloc[0]
            pred = results["course_id"].iloc[0]
            order = np.array(pred_ranking).argsort()[::-1]
            pred = np.array(pred)[order]

        pred = [self.id_to_povinn.get(cid) for cid in pred if cid in self.id_to_povinn]

        return pred

    def graph_with_user(self, user: User) -> tuple[int, HeteroData]:
        uid = self.train_data.user.shape[0]
        finished = user.blueprint_to_df()
        finished["user_id"] = uid
        finished["zskr"] = None
        finished["zroc"] = None
        finished = finished.merge(
            self.train_data.povinn[["course_id", "povinn", "embed"]],
            left_on="course",
            right_on="povinn",
            how="inner",
        )

        user_df = pd.DataFrame(
            [
                [
                    uid,
                    user.id,
                    None,  # sident
                    None,  # sdruh
                    user.enrollment_year,
                    None,  # sobor
                    None,  # sobor_nazev
                    user.degree_plan,
                    finished["embed"].values.mean(axis=0)
                    if not finished.empty
                    else self.train_data.no_history_user_embed,
                ]
            ],
            index=[uid],
            columns=self.train_data.user.columns,
        )
        user_features = get_user_features(user_df, self.train_data.user)
        edge_index = torch.stack(
            [
                torch.tensor(finished["user_id"].values, dtype=torch.long),
                torch.tensor(finished["course_id"].values, dtype=torch.long),
            ]
        )
        graph = HeteroData()
        graph["user"].x = torch.cat([self.train_graph["user"].x, user_features], dim=0)
        graph["course"].x = self.train_graph["course"].x
        graph["user", "finished", "course"].edge_index = torch.stack(
            [
                torch.cat(
                    [
                        self.train_graph["user", "finished", "course"].edge_index[0],
                        edge_index[0],
                    ],
                    dim=0,
                ),
                torch.cat(
                    [
                        self.train_graph["user", "finished", "course"].edge_index[1],
                        edge_index[1],
                    ],
                    dim=0,
                ),
            ]
        )
        graph["user", "finished", "course"].edge_label_index = torch.stack(
            [
                torch.full([self.train_data.povinn["course_id"].shape[0]], uid),
                torch.tensor(self.train_data.povinn["course_id"].values),
            ]
        )
        graph["course", "rev_finished", "user"].edge_index = torch.stack(
            [
                torch.cat(
                    [
                        self.train_graph["course", "rev_finished", "user"].edge_index[
                            0
                        ],
                        edge_index[1],
                    ],
                    dim=0,
                ),
                torch.cat(
                    [
                        self.train_graph["course", "rev_finished", "user"].edge_index[
                            1
                        ],
                        edge_index[0],
                    ],
                    dim=0,
                ),
            ]
        )
        return uid, graph

    def set_train_params(self, lr=1e-2, epochs=30, loss_fn=bpr_loss):
        self.loss_fn = loss_fn
        self.lr = lr
        self.epochs = epochs


class GNNEncoder(torch.nn.Module):
    def __init__(self, hidden_channels, out_channels, metadata):
        super().__init__()
        # metadata[1] contains edge types as tuples: (src, rel, dst)
        convs1 = {}
        convs2 = {}
        for src, rel, dst in metadata[1]:
            # Create a SAGEConv for each relation for two layers
            convs1[(src, rel, dst)] = SAGEConv((-1, -1), hidden_channels)
            convs2[(src, rel, dst)] = SAGEConv((-1, -1), out_channels)

        self.conv1 = HeteroConv(convs1, aggr="sum")
        self.conv2 = HeteroConv(convs2, aggr="sum")

    def forward(self, x_dict, edge_index_dict):
        # conv1 -> relu per node type -> conv2
        x_dict = self.conv1(x_dict, edge_index_dict)
        # apply relu to every node type tensor
        x_dict = {k: v.relu() for k, v in x_dict.items()}
        x_dict = self.conv2(x_dict, edge_index_dict)
        return x_dict


class EdgeDecoder(torch.nn.Module):
    def __init__(self, hidden_channels):
        super().__init__()
        self.lin1 = torch.nn.Linear(3 * hidden_channels, hidden_channels)
        self.dropout = torch.nn.Dropout(p=0.3)
        self.lin2 = torch.nn.Linear(hidden_channels, 1)

    def forward(self, z_dict, edge_label_index):
        row, col = edge_label_index
        z = torch.cat(
            [
                z_dict["user"][row],
                z_dict["course"][col],
                z_dict["user"][row] * z_dict["course"][col],
            ],
            dim=-1,
        )

        z = self.lin1(z)
        z = z.relu()
        z = self.lin2(z)
        return z.view(-1)


class Model(torch.nn.Module):
    def __init__(self, data, hidden_channels, out_channels):
        super().__init__()
        # self.node_emb = torch.nn.ModuleDict(
        #     {
        #         "user": torch.nn.Embedding(
        #             num_embeddings=data["user"].num_nodes, embedding_dim=hidden_channels
        #         ),
        #         "course": torch.nn.Embedding(
        #             num_embeddings=data["course"].num_nodes,
        #             embedding_dim=hidden_channels,
        #         ),
        #     }
        # )
        self.user_proj = torch.nn.Sequential(
            torch.nn.Linear(data["user"].x.size(1), hidden_channels),
            torch.nn.ReLU(),
            torch.nn.LayerNorm(hidden_channels),
            torch.nn.Dropout(p=0.2),
        )
        self.course_proj = torch.nn.Sequential(
            torch.nn.Linear(data["course"].x.size(1), hidden_channels),
            torch.nn.ReLU(),
            torch.nn.LayerNorm(hidden_channels),
            torch.nn.Dropout(p=0.2),
        )
        self.encoder = GNNEncoder(hidden_channels, out_channels, data.metadata())
        self.decoder = EdgeDecoder(out_channels)

    def forward(self, x_dict, edge_index_dict, edge_label_index):
        # x_dict = {
        #     "user": self.node_emb["user"].weight,
        #     "course": self.node_emb["course"].weight,
        # }
        #

        user_emb = self.user_proj(x_dict["user"])
        course_emb = self.course_proj(x_dict["course"])

        user_emb = torch.nn.functional.normalize(user_emb, dim=-1)
        course_emb = torch.nn.functional.normalize(course_emb, dim=-1)

        x_proj = {"user": user_emb, "course": course_emb}
        z_dict = self.encoder(x_proj, edge_index_dict)
        return self.decoder(z_dict, edge_label_index)


def train_model(model, optimizer, train_data, val_data, loss_fn, epochs):
    def train():
        model.train()
        optimizer.zero_grad()
        pred = model(
            train_data.x_dict,
            train_data.edge_index_dict,
            train_data["user", "course"].edge_label_index,
        )
        target = train_data["user", "course"].edge_label.float()
        loss = loss_fn(pred, target)
        loss.backward()
        optimizer.step()
        return float(loss.detach())

    @torch.no_grad()
    def test(data):
        data = data.to(device)
        model.eval()
        pred = model(
            data.x_dict, data.edge_index_dict, data["user", "course"].edge_label_index
        )
        target = data["user", "course"].edge_label.float()
        loss = loss_fn(pred, target)
        return float(loss)

    for epoch in range(1, epochs):
        train_data = train_data.to(device)
        loss = train()
        train_loss = test(train_data)
        val_loss = test(val_data)

        print(
            f"Epoch: {epoch:03d}, Loss: {loss:.4f}, Train: {train_loss:.4f}, "
            f"Val: {val_loss:.4f}"
        )


def negative_split(
    train, val, test, all_interaction, train_ratio=1, val_ratio=1, test_ratio=-1
):
    def negative(interaction, all_interaction, ratio):
        course = pd.Series(interaction["course_id"].unique())
        if ratio > 0:
            user_interaction_count = (
                interaction.groupby("user_id")["course_id"].count() * ratio
            ).round()
        else:
            user_interaction_count = pd.Series(
                -1, index=pd.Index(interaction["user_id"].unique(), name="user_id")
            )
        user_interaction_count = (
            user_interaction_count.rename("count").reset_index().astype(int)
        )

        negative_per_user = pd.Series(
            user_interaction_count.apply(
                lambda x: (
                    course[
                        ~course.isin(
                            all_interaction[all_interaction["user_id"] == x["user_id"]][
                                "course_id"
                            ]
                        )
                    ]
                    .sample(
                        x["count"] if x["count"] > 0 else None,
                        frac=(None if x["count"] > 0 else 1),
                        random_state=RND_STATE,
                    )
                    .tolist()
                ),
                axis=1,
            ).values,
            index=user_interaction_count["user_id"],
        )
        negative = (
            negative_per_user.explode()
            .reset_index()
            .rename(columns={"index": "user_id", 0: "course_id"})
        )
        negative = negative.dropna()
        negative["zskr"] = 0
        negative["course_id"] = negative["course_id"].astype(int)
        return negative

    neg_train = negative(train, all_interaction, train_ratio)
    neg_val = negative(val, all_interaction, val_ratio)
    neg_test = None
    if test is not None:
        neg_test = negative(test, all_interaction, test_ratio)

    return neg_train, neg_val, neg_test


def get_user_features(user: pd.DataFrame, all_users: pd.DataFrame) -> torch.Tensor:
    study_type = pd.crosstab(user["user_id"], user["sdruh"])
    study_type = study_type.reindex(
        index=user["user_id"],
        columns=all_users["sdruh"].sort_values().unique(),
        fill_value=0,
    )
    study_type = torch.from_numpy(study_type.values).to(torch.float)
    field = pd.crosstab(user["user_id"], user["sobor"])
    field = field.reindex(
        index=user["user_id"],
        columns=all_users["sobor"].sort_values().unique(),
        fill_value=0,
    )
    field = torch.from_numpy(field.values).to(torch.float)
    embed = np.stack(user["embed"], axis=0)
    embed = embed.astype(np.float32)
    embed = torch.from_numpy(embed)
    user_features = torch.cat([embed, study_type, field], dim=-1)
    return user_features


def get_course_features(course: pd.DataFrame) -> torch.Tensor:
    # department = pd.get_dummies(course["pgarant"])
    # department = torch.from_numpy(department.values).to(torch.float)
    embed = np.stack(course["embed"], axis=0)
    embed = embed.astype(np.float32)
    embed = torch.from_numpy(embed)
    course_features = torch.cat([embed], dim=-1)

    return course_features


def graph_data(
    user, course, pos_train, pos_val, pos_test, neg_train, neg_val, neg_test
):
    def index_label(pos, neg):
        df = pd.concat([pos, neg])
        index = torch.stack(
            [
                torch.tensor(df["user_id"].values, dtype=torch.long),
                torch.tensor(df["course_id"].values, dtype=torch.long),
            ]
        )
        label = torch.cat(
            [
                torch.ones(pos.shape[0], dtype=torch.long),
                torch.zeros(neg.shape[0], dtype=torch.long),
            ]
        )
        return index, label

    def hetero_data_builder(user_features, course_features, edge_index):
        def hetero_data(edge_label_index, edge_label):
            data = HeteroData()
            data["user"].x = user_features
            data["course"].x = course_features
            data["user", "finished", "course"].edge_index = edge_index
            data["user", "finished", "course"].edge_label = edge_label
            data["user", "finished", "course"].edge_label_index = edge_label_index
            data = T.ToUndirected()(data)
            del data["course", "rev_finished", "user"].edge_label
            return data

        return hetero_data

    user_features = get_user_features(user, user)
    course_features = get_course_features(course)

    # Edge
    edge_index, _ = index_label(pos_train, pd.DataFrame())
    train_index, train_label = index_label(pos_train, neg_train)
    val_index, val_label = index_label(pos_val, neg_val)

    hdb = hetero_data_builder(user_features, course_features, edge_index)
    train = hdb(train_index, train_label)
    val = hdb(val_index, val_label)

    test = None
    if test is not None:
        test_index, test_label = index_label(pos_test, neg_test)
        test = hdb(test_index, test_label)

    return train, val, test
