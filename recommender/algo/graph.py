import numpy as np
import pandas as pd
import torch
import torch.nn.functional as F
import torch_geometric.transforms as T
from torch_geometric.data import HeteroData
from torch_geometric.nn import SAGEConv, to_hetero

from algo.train import TrainData
from user import User

RND_STATE = 42
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


class GraphSage(TrainData):
    def fit(self):
        LOSS_FN = bpr_loss
        VAL_RATIO = 0.2
        lr = 1e-2
        epochs = 1

        user, finished, povinn = self.dataset()
        train, val, test = self.split(finished, VAL_RATIO)
        neg_train, neg_val, neg_test = negative_split(
            train, val, test, finished, val_ratio=-1, train_ratio=1
        )
        train, val, test = graph_data(
            user, povinn, train, val, test, neg_train, neg_val, neg_test
        )

        self.model = Model(train, hidden_channels=32, out_channels=32).to(device)

        optimizer = torch.optim.AdamW(self.model.parameters(), lr=lr)
        train_model(self.model, optimizer, train, val, LOSS_FN, epochs)

    def recommend(self, user: User, limit: int) -> list[str]:
        # TODO: implement
        results = predict(model, eval_data, LOSS_FN)

        results = results.groupby("user_id").agg({"pred": list, "target": list})
        results["target"] = results["target"].apply(np.array)
        results["pred"] = results["pred"].apply(np.array)
        results["sort_ids"] = results["pred"].apply(lambda x: x.argsort())
        results["target"] = results.apply(
            lambda x: x["target"][x["sort_ids"][::-1]], axis=1
        )
        results["pred"] = results.apply(
            lambda x: x["pred"][x["sort_ids"][::-1]], axis=1
        )
        results = results.drop(columns=["sort_ids"])
        results = results.reset_index()


def predict(model, test_data, loss_fn):
    with torch.no_grad():
        test_data = test_data.to(device)
        pred = model(
            test_data.x_dict,
            test_data.edge_index_dict,
            test_data["user", "course"].edge_label_index,
        )
        target = test_data["user", "course"].edge_label.float()
        bce = loss_fn(pred, target)
        print(f"Eval Loss: {bce:.4f}")

    user_id = test_data["user", "course"].edge_label_index[0].cpu().numpy()
    course_id = test_data["user", "course"].edge_label_index[1].cpu().numpy()

    results = pd.DataFrame(
        {
            "user_id": user_id,
            "course_id": course_id,
            "pred": pred.sigmoid().cpu().numpy(),
            "target": target.cpu().numpy(),
        }
    )
    return results


class GNNEncoder(torch.nn.Module):
    def __init__(self, hidden_channels, out_channels):
        super().__init__()
        self.conv1 = SAGEConv((-1, -1), hidden_channels)
        self.conv2 = SAGEConv((-1, -1), out_channels)

    def forward(self, x, edge_index):
        x = self.conv1(x, edge_index)
        x = x.relu()
        x = self.conv2(x, edge_index)
        return x


class EdgeDecoder(torch.nn.Module):
    def __init__(self, hidden_channels):
        super().__init__()
        self.lin1 = torch.nn.Linear(2 * hidden_channels, hidden_channels)
        self.dropout = torch.nn.Dropout(p=0.1)
        self.lin2 = torch.nn.Linear(hidden_channels, 1)

    def forward(self, z_dict, edge_label_index):
        row, col = edge_label_index
        z = torch.cat([z_dict["user"][row], z_dict["course"][col]], dim=-1)

        z = self.lin1(z)
        z = z.relu()
        z = self.lin2(z)
        return z.view(-1)


class Model(torch.nn.Module):
    def __init__(self, data, hidden_channels, out_channels):
        super().__init__()
        self.encoder = GNNEncoder(hidden_channels, out_channels)
        self.encoder = to_hetero(self.encoder, data.metadata(), aggr="sum")
        self.decoder = EdgeDecoder(out_channels)

    def forward(self, x_dict, edge_index_dict, edge_label_index):
        z_dict = self.encoder(x_dict, edge_index_dict)
        return self.decoder(z_dict, edge_label_index)


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
    torch.save(model.state_dict(), "model.pth")


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
            user_interaction_count = interaction.groupby("user_id").apply(lambda x: -1)
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
    neg_test = negative(test, all_interaction, test_ratio)

    return neg_train, neg_val, neg_test


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

    # User
    study_type = pd.get_dummies(user["sdruh"])
    study_type = torch.from_numpy(study_type.values).to(torch.float)
    field = pd.get_dummies(user["sobor"])
    field = torch.from_numpy(field.values).to(torch.float)
    # field_embed = torch.tensor(user["sobor_embed"].tolist())
    user_features = torch.cat([study_type, field], dim=-1)

    # Course
    department = pd.get_dummies(course["pgarant"])
    department = torch.from_numpy(department.values).to(torch.float)
    # name = course["pnazev_embed"]
    # name = torch.tensor(name.tolist())
    # name_embed = torch.tensor(course["pnazev_embed"].tolist())
    course_features = torch.cat([department], dim=-1)

    # Edge
    edge_index, _ = index_label(pos_train, pd.DataFrame())
    train_index, train_label = index_label(pos_train, neg_train)
    val_index, val_label = index_label(pos_val, neg_val)
    test_index, test_label = index_label(pos_test, neg_test)

    hdb = hetero_data_builder(user_features, course_features, edge_index)
    train = hdb(train_index, train_label)
    val = hdb(val_index, val_label)
    test = hdb(test_index, test_label)

    return train, val, test
