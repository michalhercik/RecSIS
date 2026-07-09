import argparse
import os
import sys

import numpy as np
import pandas as pd
import torch
import torch.nn.functional as F
import torch_geometric.transforms as T
from embedder import sbert_embed
from retrieve import user_interaction_povinn
from torch_geometric.data import HeteroData
from torch_geometric.nn import HeteroConv, SAGEConv

sys.path.insert(0, "..")
from data_repository import DataRepository

RND_STATE = 42

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


def main(args):
    LOSS_FN = bpr_loss
    VAL_RATIO = 0.2

    user, finished, povinn = dataset(args.data)
    train, val, test = split(finished, VAL_RATIO)
    neg_train, neg_val, neg_test = negative_split(
        train, val, test, finished, val_ratio=-1, train_ratio=1
    )

    user = user.merge(
        train.merge(povinn[["course_id", "embed"]], on="course_id")
        .groupby("user_id")
        .agg(
            {
                "user_id": "first",
                "course_id": set,
                "embed": lambda x: np.mean(x.values, axis=0),
            }
        )[["user_id", "embed"]],
        left_on="user_id",
        right_index=True,
        how="left",
    )
    empty = list(sbert_embed([""]))[0]

    def _fill_embed(x):
        na = pd.isna(x)
        if isinstance(na, (np.ndarray, list, pd.Series)):
            na = bool(np.all(na))
        if na:
            return empty
        return x

    user["embed"] = user["embed"].apply(_fill_embed)

    train, val, test = graph_data(
        user, povinn, train, val, test, neg_train, neg_val, neg_test
    )

    print(train)
    print(val if not args.eval else test)

    model = Model(train, hidden_channels=64, out_channels=32).to(device)
    print(model)

    if os.path.exists("model.pth") and not args.train:
        model.load_state_dict(torch.load("model.pth", weights_only=True))
    else:
        optimizer = torch.optim.AdamW(model.parameters(), lr=args.lr)
        train_model(model, optimizer, train, val, LOSS_FN, args.epochs)

    eval_data = val if not args.eval else test
    results = predict(model, eval_data, LOSS_FN)

    results = results.groupby("user_id").agg({"pred": list, "target": list})
    results["target"] = results["target"].apply(np.array)
    results["pred"] = results["pred"].apply(np.array)
    results["sort_ids"] = results["pred"].apply(lambda x: x.argsort())
    results["target"] = results.apply(
        lambda x: x["target"][x["sort_ids"][::-1]], axis=1
    )
    results["pred"] = results.apply(lambda x: x["pred"][x["sort_ids"][::-1]], axis=1)
    results = results.drop(columns=["sort_ids"])
    results = results.reset_index()

    results_description = eval(user, results)
    print(results_description)


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
        self.node_emb = torch.nn.ModuleDict(
            {
                "user": torch.nn.Embedding(
                    num_embeddings=data["user"].num_nodes, embedding_dim=hidden_channels
                ),
                "course": torch.nn.Embedding(
                    num_embeddings=data["course"].num_nodes,
                    embedding_dim=hidden_channels,
                ),
            }
        )
        self.encoder = GNNEncoder(hidden_channels, out_channels, data.metadata())
        self.decoder = EdgeDecoder(out_channels)

    def forward(self, x_dict, edge_index_dict, edge_label_index):
        x_dict = {
            "user": self.node_emb["user"].weight,
            "course": self.node_emb["course"].weight,
        }
        z_dict = self.encoder(x_dict, edge_index_dict)
        return self.decoder(z_dict, edge_label_index)


def recall(pred, target, k=-1):
    assert pred.shape == target.shape
    target_k = target
    if k > 0:
        assert target.shape[0] >= k
        target_k = target[:k]
    return target_k.sum() / target.sum()


def precision(pred, target, k=-1):
    assert pred.shape == target.shape
    if k > 0:
        assert pred.shape[0] >= k
        pred = pred[:k]
        target = target[:k]
    pred = pred.round()
    return pred[pred == target].sum() / target.shape[0]


def map(pred, target, k=-1):
    if k > 0:
        assert target.shape[0] >= k
        target = target[:k]
    cumsum = np.cumsum(target)
    hit = np.where(target == 1)[0]
    return (cumsum[hit] / (hit + 1)).mean() if hit.size > 0 else 0


def eval(user: pd.DataFrame, results: pd.DataFrame):
    metrics = []
    # for k in [5, 10, 20, 50]:
    for k in [5, 24, 50]:
        r = results.apply(lambda x: recall(x["pred"], x["target"], k), axis=1)
        p = results.apply(lambda x: precision(x["pred"], x["target"], k), axis=1)
        m = results.apply(lambda x: map(x["pred"], x["target"], k), axis=1)
        metrics.append(
            pd.DataFrame(
                {
                    "user_id": results["user_id"],
                    "k": k,
                    "Recall": r,
                    "Precision": p,
                    "AP": m,
                }
            )
        )

    results = pd.concat(metrics)
    results = results.merge(
        user[["user_id", "soident", "sdruh", "sobor"]], on="user_id"
    )

    # describe = results.groupby(["k", "sdruh"])[["Recall", "Precision", "AP"]].describe()
    describe = results.groupby(["k", "sdruh"])[["Recall", "AP"]].describe()
    print(describe)
    describe = results.groupby(["k"])[["Recall", "AP"]].describe()
    describe_all = (
        describe.mean().to_frame().T.set_index([pd.Index(["-"]), pd.Index(["-"])])
    )
    describe = pd.concat([describe, describe_all]).round(4)
    return describe


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


def bce_loss(pred, target):
    return F.binary_cross_entropy_with_logits(pred, target)


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


def dataset(force_load=True):
    def from_db():
        user, interaction, povinn = user_interaction_povinn()

        user = user.reset_index().rename(columns={"index": "user_id"})

        povinn = povinn.reset_index().rename(columns={"index": "course_id"})
        pamela = DataRepository().pamela
        pamela = pamela[
            (pamela["jazyk"] == "ENG") & (pamela["typ"].isin(["A", "S"]))
        ].pivot_table(index="povinn", columns="typ", values="memo", aggfunc="first")
        povinn = povinn.merge(pamela, on="povinn", how="left")
        embed_src = povinn.apply(
            lambda x: f"{x['panazev']}: {x['A']}\n{x['S']}", axis=1
        )
        povinn["embed"] = list(sbert_embed(embed_src))

        interaction = interaction.merge(user[["sident", "user_id"]], on="sident")
        interaction = interaction.merge(povinn[["povinn", "course_id"]], on="povinn")
        interaction = interaction[["user_id", "course_id", "zskr"]]

        return user, interaction, povinn

    # look for separated files: metadata + embeddings
    meta_path = "povinn_meta.pkl"
    emb_path = "povinn_emb.npy"

    if (
        os.path.exists("user.csv")
        and os.path.exists(meta_path)
        and os.path.exists(emb_path)
        and os.path.exists("interaction.csv")
        and not force_load
    ):
        user = pd.read_csv("user.csv")
        povinn = pd.read_pickle(meta_path)
        # load embeddings matrix and attach back as rows (numpy arrays)
        embs = np.load(emb_path)
        # make sure lengths match
        if len(embs) != len(povinn):
            raise RuntimeError(
                "Embeddings file length does not match povinn metadata length"
            )
        povinn["embed"] = list(embs)
        interaction = pd.read_csv("interaction.csv")
    else:
        user, interaction, povinn = from_db()
        user.to_csv("user.csv", index=False)

        # make a serializable copy: convert torch tensors to numpy if needed
        serializable_povinn = povinn.copy()
        serializable_povinn["embed"] = serializable_povinn["embed"].apply(
            lambda x: x.numpy() if hasattr(x, "numpy") else np.asarray(x)
        )

        # stack embeddings into a single 2D array and save as .npy (float32)
        embs = np.vstack(serializable_povinn["embed"].values).astype(np.float32)
        np.save(emb_path, embs)

        # save metadata without the embedding column
        meta = serializable_povinn.drop(columns=["embed"])
        meta.to_pickle(meta_path)

        interaction.to_csv("interaction.csv", index=False)

        # restore original povinn (with embeddings) to return
        povinn = meta.copy()
        povinn["embed"] = list(embs)

    return user, interaction, povinn


def split(interaction, val_ratio, split_year=2024):
    # Train data are all interactions before split_year
    train = interaction[interaction["zskr"] < split_year]

    # Test data are all interactions after split_year (including split_year)
    year_bitmap = interaction["zskr"] >= split_year

    # Split test data using val_ratio into validation and test sets by user_id randomly
    test_user_id = interaction[year_bitmap]["user_id"].drop_duplicates()
    val_user_id = test_user_id.sample(frac=val_ratio, random_state=RND_STATE)
    val_user_bitmap = interaction["user_id"].isin(val_user_id)
    test_user_bitmap = interaction["user_id"].isin(test_user_id.drop(val_user_id.index))

    val = interaction[year_bitmap & val_user_bitmap]
    test = interaction[year_bitmap & test_user_bitmap]
    # val = interaction[test_bitmap].sample(frac=val_ratio)
    # test = interaction[test_bitmap].drop(val.index)

    return train, val, test


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
    print(user.shape)
    print(user["embed"])
    emebed = torch.tensor(user["embed"].tolist())
    user_features = torch.cat([emebed, study_type, field], dim=-1)

    # Course
    department = pd.get_dummies(course["pgarant"])
    department = torch.from_numpy(department.values).to(torch.float)
    embed = torch.tensor(course["embed"].tolist())
    course_features = torch.cat([embed], dim=-1)

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


def parser():
    parser = argparse.ArgumentParser(
        prog="RECSIS Graph Recommender",
        description="Train and evaluate the model",
    )
    parser.add_argument(
        "-t",
        "--train",
        action="store_true",
        default=False,
        help="Train the model even if it already exists",
    )
    parser.add_argument(
        "-e",
        "--epochs",
        type=int,
        default=10,
        help="Number of epochs to train the model",
    )
    parser.add_argument(
        "-lr",
        "--learning-rate",
        type=float,
        default=1e-2,
        dest="lr",
        help="Learning rate for the optimizer",
    )
    parser.add_argument(
        "-v",
        "--eval",
        action="store_true",
        default=False,
        help="Evaluate the model on the test set",
    )
    parser.add_argument(
        "-d",
        "--data",
        action="store_true",
        default=False,
        help="Load data from database",
    )
    return parser


if __name__ == "__main__":
    args = parser().parse_args()
    main(args)

