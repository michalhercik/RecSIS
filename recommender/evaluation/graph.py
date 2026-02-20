import argparse
import os

import numpy as np
import pandas as pd
import torch
import torch.nn.functional as F
import torch_geometric.transforms as T
from retrieve import user_interaction_povinn
from torch_geometric.data import HeteroData
from torch_geometric.nn import SAGEConv, to_hetero

# python .\graph.py -t -e 300 -lr 1e-2
# Eval Loss: 0.4599
#      Recall                                                      Precision                                                          AP
#       count    mean     std  min     25%     50%     75%     max     count    mean     std  min     25%     50%     75%     max  count    mean     std  min     25%     50%     75%     max
# 5  B  172.0  0.1727  0.1419  0.0  0.0692  0.1818  0.2500  1.0000     172.0  0.4140  0.3424  0.0  0.0000  0.4000  0.8000  1.0000  172.0  0.5791  0.3948  0.0  0.2000  0.6896  0.9500  1.0000
#    N   65.0  0.0528  0.0994  0.0  0.0000  0.0000  0.0667  0.5714      65.0  0.1046  0.1662  0.0  0.0000  0.0000  0.2000  0.8000   65.0  0.1654  0.2449  0.0  0.0000  0.0000  0.2500  1.0000
# 10 B  172.0  0.3311  0.2049  0.0  0.1818  0.3333  0.4706  1.0000     172.0  0.3436  0.3153  0.0  0.1000  0.2000  0.6250  1.0000  172.0  0.5549  0.3349  0.0  0.2482  0.6222  0.8553  1.0000
#    N   65.0  0.1014  0.1473  0.0  0.0000  0.0526  0.1667  0.8000      65.0  0.0769  0.1260  0.0  0.0000  0.0000  0.1000  0.5000   65.0  0.1815  0.2212  0.0  0.0000  0.1000  0.3028  1.0000
# 20 B  172.0  0.5174  0.2420  0.0  0.3333  0.5714  0.6667  1.0000     172.0  0.2142  0.2181  0.0  0.0500  0.1000  0.4000  0.6500  172.0  0.5057  0.2942  0.0  0.2475  0.5325  0.8191  0.9924
#    N   65.0  0.1745  0.1756  0.0  0.0000  0.1739  0.2667  0.8000      65.0  0.0454  0.0722  0.0  0.0000  0.0000  0.0500  0.3000   65.0  0.1814  0.1927  0.0  0.0000  0.1026  0.3026  0.8769
# 50 B  172.0  0.7670  0.2281  0.0  0.6250  0.8209  0.9506  1.0000     172.0  0.1020  0.1144  0.0  0.0200  0.0500  0.1600  0.3800  172.0  0.4416  0.2543  0.0  0.2286  0.4254  0.7091  0.8936
#    N   65.0  0.4493  0.2325  0.0  0.3333  0.4167  0.6111  1.0000      65.0  0.0188  0.0293  0.0  0.0000  0.0000  0.0200  0.1200   65.0  0.1748  0.1558  0.0  0.0542  0.1102  0.2760  0.8769
# -  -  118.5  0.3208  0.1840  0.0  0.1928  0.3188  0.4311  0.8964     118.5  0.1649  0.1730  0.0  0.0213  0.0938  0.2944  0.5938  118.5  0.3480  0.2616  0.0  0.1223  0.3228  0.5581  0.9550

# df = load(OUT)
# me = df[
#     df["relevant"].apply(
#         lambda x: set(
#             ["NMAT100", "NPRG051", "NSWI153", "NSWI098", "NSWI041", "NSWI154"]
#         )
#         <= set(x)
#     )
#     & (df["zroc"] == 1)
#     & (df["zsem"] == 1)
# ]
# print(me)

RND_STATE = 42

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


def main(args):
    # LOSS_FN = margin_ranking_loss
    LOSS_FN = bpr_loss
    VAL_RATIO = 0.2

    user, finished, povinn = dataset(args.data)
    train, val, test = split(finished, VAL_RATIO)
    neg_train, neg_val, neg_test = negative_split(
        train, val, test, finished, val_ratio=-1, train_ratio=1
    )
    train, val, test = graph_data(
        user, povinn, train, val, test, neg_train, neg_val, neg_test
    )
    print(train)
    print(val if not args.eval else test)

    model = Model(train, hidden_channels=32, out_channels=32).to(device)
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
    metrics = []
    for k in [5, 10, 20, 50]:
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

    describe = results.groupby(["k", "sdruh"])[["Recall", "Precision", "AP"]].describe()
    describe_all = (
        describe.mean().to_frame().T.set_index([pd.Index(["-"]), pd.Index(["-"])])
    )
    describe = pd.concat([describe, describe_all]).round(4)
    print(describe)


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
    for k in [5, 10, 20, 50]:
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

    describe = results.groupby(["k", "sdruh"])[["Recall", "Precision", "AP"]].describe()
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


def margin_ranking_loss(
    pred,
    target,
    num_negatives=10,
    margin=1.0,
):
    """
    pred   : Tensor [N]
    target : Tensor [N] ∈ {0,1}
    """

    pred = pred.view(-1)
    target = target.view(-1)

    pos_pred = pred[target == 1]
    neg_pred = pred[target == 0]

    if len(pos_pred) == 0 or len(neg_pred) == 0:
        return torch.tensor(0.0, device=pred.device, requires_grad=True)

    # sample negatives for each positive
    neg_idx = torch.randint(
        0,
        neg_pred.size(0),
        (pos_pred.size(0), num_negatives),
        device=pred.device,
    )

    sampled_neg = neg_pred[neg_idx]  # [P, K]

    pos_pred = pos_pred.unsqueeze(1).expand_as(sampled_neg)

    y = torch.ones_like(pos_pred)

    return F.margin_ranking_loss(
        pos_pred.reshape(-1),
        sampled_neg.reshape(-1),
        y.reshape(-1),
        margin=margin,
        reduction="mean",
    )


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
        # user["sobor_embed"] = list(sbert_embed(user["sobor_nazev"]))
        povinn = povinn.reset_index().rename(columns={"index": "course_id"})
        # povinn["pnazev_embed"] = list(sbert_embed(povinn["pnazev"]))

        interaction = interaction.merge(user[["sident", "user_id"]], on="sident")
        interaction = interaction.merge(povinn[["povinn", "course_id"]], on="povinn")
        interaction = interaction[["user_id", "course_id", "zskr"]]

        return user, interaction, povinn

    if (
        os.path.exists("user.csv")
        and os.path.exists("povinn.pkl")
        and os.path.exists("interaction.csv")
        and not force_load
    ):
        user = pd.read_csv("user.csv")
        povinn = pd.read_pickle("povinn.pkl")
        # povinn["pnazev_embed"] = povinn["pnazev_embed"].apply(torch.tensor)
        interaction = pd.read_csv("interaction.csv")
    else:
        user, interaction, povinn = from_db()
        user.to_csv("user.csv")
        serializable_povinn = povinn
        # serializable_povinn["pnazev_embed"] = serializable_povinn["pnazev_embed"].apply(
        #     lambda x: x.numpy()
        # )
        serializable_povinn.to_pickle("povinn.pkl")
        interaction.to_csv("interaction.csv")
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
                lambda x: course[
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
                .tolist(),
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


# def dataset():
#     df = load(OUT)
#     df["interacts"] = df.apply(lambda x: x["finished"] + x["relevant"], axis=1)

#     df = (
#         df[["soident", "sdruh", "sobor", "interacts"]]
#         .explode("interacts")
#         .groupby(["soident", "sdruh", "sobor"])
#         .agg(list)
#     )
#     df = df.rename(columns={"interacts": "finished"})
#     df = df.reset_index()
#     df["finished"] = df["finished"].apply(set)
#     df = df.reset_index().rename(columns={"index": "user_id"})

#     finished = (
#         df[["user_id", "finished"]]
#         .explode("finished")
#         .rename(columns={"finished": "course"})
#     )
#     povinn2id = (
#         pd.DataFrame(finished["course"].unique(), columns=["course"])
#         .reset_index()
#         .rename(columns={"index": "course_id"})
#     )
#     povinn = DataRepository().povinn
#     povinn = povinn[
#         ["povinn", "pnazev", "panazev", "pfakulta", "pgarant", "pvyucovan", "vsemzac"]
#     ].rename(columns={"povinn": "course"})
#     povinn = povinn2id.merge(povinn, on="course")

#     finished = finished.merge(povinn2id, on="course")

#     return df, finished, povinn


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

# Proper data split, epoch=200, lr=0.01, val_ratio=0.2, loss=bpr
#
# Eval Loss: 0.5240
#         Recall                                                   Precision                                                    Average precision
#          count    mean     std  min     25%     50%     75%  max     count    mean     std  min     25%     50%     75%   max             count    mean     std  min     25%     50%     75%  max
# 5   B    754.0  0.2726  0.3051  0.0  0.0000  0.2500  0.5000  1.0     754.0  0.1581  0.1744  0.0  0.0000  0.2000  0.2000  0.80             754.0  0.2910  0.3347  0.0  0.0000  0.2500  0.5000  1.0
#     N    255.0  0.0627  0.1911  0.0  0.0000  0.0000  0.0000  1.0     255.0  0.0282  0.0762  0.0  0.0000  0.0000  0.0000  0.40             255.0  0.0650  0.2007  0.0  0.0000  0.0000  0.0000  1.0
# 10  B    754.0  0.3868  0.3498  0.0  0.0000  0.3333  0.6667  1.0     754.0  0.1146  0.1160  0.0  0.0000  0.1000  0.2000  0.50             754.0  0.2938  0.3050  0.0  0.0000  0.2500  0.5000  1.0
#     N    255.0  0.1336  0.2886  0.0  0.0000  0.0000  0.0000  1.0     255.0  0.0286  0.0575  0.0  0.0000  0.0000  0.0000  0.30             255.0  0.0756  0.1905  0.0  0.0000  0.0000  0.0000  1.0
# 20  B    754.0  0.5312  0.3696  0.0  0.2500  0.5000  1.0000  1.0     754.0  0.0789  0.0689  0.0  0.0500  0.0500  0.1000  0.35             754.0  0.2781  0.2720  0.0  0.0556  0.2190  0.4496  1.0
#     N    255.0  0.2522  0.3491  0.0  0.0000  0.0000  0.5000  1.0     255.0  0.0284  0.0375  0.0  0.0000  0.0000  0.0500  0.15             255.0  0.0856  0.1713  0.0  0.0000  0.0000  0.0909  1.0
# 50  B    754.0  0.7534  0.3380  0.0  0.5000  1.0000  1.0000  1.0     754.0  0.0437  0.0315  0.0  0.0200  0.0400  0.0600  0.18             754.0  0.2454  0.2336  0.0  0.0594  0.1870  0.3603  1.0
#     N    255.0  0.4819  0.4002  0.0  0.0000  0.5000  1.0000  1.0     255.0  0.0221  0.0204  0.0  0.0000  0.0200  0.0400  0.12             255.0  0.0884  0.1542  0.0  0.0000  0.0384  0.0909  1.0
# ALL ALL  504.5  0.3593  0.3239  0.0  0.0938  0.3229  0.5833  1.0     504.5  0.0628  0.0728  0.0  0.0088  0.0512  0.0813  0.35             504.5  0.1779  0.2327  0.0  0.0144  0.1181  0.2490  1.0
#

# Proper data split, epoch=20, lr=0.01, val_ratio=0.2, loss=bpr
#
# Eval Loss: 0.5404
#      Recall                                                   Precision                                                  Average precision
#       count    mean     std  min  25%     50%     75%     max     count    mean     std  min  25%     50%    75%     max             count    mean     std  min  25%     50%     75%     max
# 5  B  754.0  0.0561  0.1296  0.0  0.0  0.0000  0.0000  1.0000     754.0  0.0411  0.0890  0.0  0.0  0.0000  0.000  0.4000             754.0  0.0885  0.1987  0.0  0.0  0.0000  0.0000  1.0000
#    N  248.0  0.0012  0.0142  0.0  0.0  0.0000  0.0000  0.2000     248.0  0.0016  0.0179  0.0  0.0  0.0000  0.000  0.2000             248.0  0.0018  0.0203  0.0  0.0  0.0000  0.0000  0.2500
# 10 B  754.0  0.1707  0.2619  0.0  0.0  0.0000  0.3333  1.0000     754.0  0.0578  0.0868  0.0  0.0  0.0000  0.100  0.4000             754.0  0.1055  0.1721  0.0  0.0  0.0000  0.1667  1.0000
#    N  248.0  0.0179  0.0770  0.0  0.0  0.0000  0.0000  0.5000     248.0  0.0069  0.0297  0.0  0.0  0.0000  0.000  0.2000             248.0  0.0087  0.0376  0.0  0.0  0.0000  0.0000  0.2917
# 20 B  754.0  0.3336  0.3563  0.0  0.0  0.3095  0.6000  1.0000     754.0  0.0526  0.0606  0.0  0.0  0.0500  0.100  0.2500             754.0  0.1069  0.1396  0.0  0.0  0.0588  0.1642  0.6806
#    N  248.0  0.0466  0.1454  0.0  0.0  0.0000  0.0000  1.0000     248.0  0.0077  0.0230  0.0  0.0  0.0000  0.000  0.1500             248.0  0.0141  0.0417  0.0  0.0  0.0000  0.0000  0.2500
# 50 B  754.0  0.5783  0.4090  0.0  0.0  0.6667  1.0000  1.0000     754.0  0.0355  0.0322  0.0  0.0  0.0200  0.060  0.2000             754.0  0.1015  0.1138  0.0  0.0  0.0639  0.1474  0.6806
#    N  248.0  0.1282  0.2437  0.0  0.0  0.0000  0.2000  1.0000     248.0  0.0082  0.0174  0.0  0.0  0.0000  0.020  0.1400             248.0  0.0192  0.0403  0.0  0.0  0.0000  0.0240  0.2001
# -  -  501.0  0.1666  0.2046  0.0  0.0  0.1220  0.2667  0.8375     501.0  0.0264  0.0446  0.0  0.0  0.0088  0.035  0.2425             501.0  0.0558  0.0955  0.0  0.0  0.0153  0.0628  0.5441


# Proper data split, epoch=20, lr=0.01, val_ratio=0.2, loss=margin_ranking_loss
#
# Eval Loss: 0.6192
#      Recall                                                   Precision                                                Average precision
#       count    mean     std  min  25%     50%     75%     max     count    mean     std  min  25%     50%    75%   max             count    mean     std  min  25%     50%     75%     max
# 5  B  761.0  0.0604  0.1474  0.0  0.0  0.0000  0.0000  1.0000     761.0  0.0405  0.0891  0.0  0.0  0.0000  0.000  0.40             761.0  0.1161  0.2846  0.0  0.0  0.0000  0.0000  1.0000
#    N  245.0  0.0050  0.0396  0.0  0.0  0.0000  0.0000  0.5000     245.0  0.0041  0.0283  0.0  0.0  0.0000  0.000  0.20             245.0  0.0043  0.0299  0.0  0.0  0.0000  0.0000  0.2500
# 10 B  761.0  0.1932  0.2746  0.0  0.0  0.0000  0.3333  1.0000     761.0  0.0624  0.0875  0.0  0.0  0.0000  0.100  0.40             761.0  0.1340  0.2328  0.0  0.0  0.0000  0.1667  1.0000
#    N  245.0  0.0253  0.1167  0.0  0.0  0.0000  0.0000  1.0000     245.0  0.0073  0.0291  0.0  0.0  0.0000  0.000  0.20             245.0  0.0112  0.0440  0.0  0.0  0.0000  0.0000  0.2667
# 20 B  761.0  0.3190  0.3572  0.0  0.0  0.2000  0.6000  1.0000     761.0  0.0503  0.0616  0.0  0.0  0.0500  0.100  0.30             761.0  0.1261  0.1914  0.0  0.0  0.0556  0.1667  1.0000
#    N  245.0  0.0387  0.1334  0.0  0.0  0.0000  0.0000  1.0000     245.0  0.0057  0.0172  0.0  0.0  0.0000  0.000  0.10             245.0  0.0139  0.0454  0.0  0.0  0.0000  0.0000  0.2667
# 50 B  761.0  0.5795  0.4143  0.0  0.0  0.6667  1.0000  1.0000     761.0  0.0351  0.0331  0.0  0.0  0.0200  0.060  0.16             761.0  0.1150  0.1507  0.0  0.0  0.0627  0.1583  1.0000
#    N  245.0  0.1283  0.2576  0.0  0.0  0.0000  0.1667  1.0000     245.0  0.0075  0.0150  0.0  0.0  0.0000  0.020  0.08             245.0  0.0171  0.0407  0.0  0.0  0.0000  0.0200  0.2500
# -  -  503.0  0.1687  0.2176  0.0  0.0  0.1083  0.2625  0.9375     503.0  0.0266  0.0451  0.0  0.0  0.0088  0.035  0.23             503.0  0.0672  0.1274  0.0  0.0  0.0148  0.0640  0.6292


# Proper data split, epoch=200, lr=0.01, val_ratio=0.2, loss=margin_ranking_loss
#
# Eval Loss: 0.4544
#      Recall                                                 Precision                                                     Average precision
#       count    mean     std  min    25%    50%     75%  max     count    mean     std  min     25%     50%    75%     max             count    mean     std  min     25%     50%     75%  max
# 5  B  744.0  0.2790  0.3235  0.0  0.000  0.200  0.5000  1.0     744.0  0.1608  0.1772  0.0  0.0000  0.2000  0.200  0.8000             744.0  0.2974  0.3307  0.0  0.0000  0.2000  0.5000  1.0
#    N  242.0  0.0844  0.2498  0.0  0.000  0.000  0.0000  1.0     242.0  0.0306  0.0886  0.0  0.0000  0.0000  0.000  0.6000             242.0  0.0601  0.1835  0.0  0.0000  0.0000  0.0000  1.0
# 10 B  744.0  0.4061  0.3673  0.0  0.000  0.400  0.6667  1.0     744.0  0.1249  0.1253  0.0  0.0000  0.1000  0.200  0.6000             744.0  0.2922  0.2948  0.0  0.0000  0.2464  0.5000  1.0
#    N  242.0  0.1257  0.2924  0.0  0.000  0.000  0.0000  1.0     242.0  0.0248  0.0573  0.0  0.0000  0.0000  0.000  0.4000             242.0  0.0700  0.1833  0.0  0.0000  0.0000  0.0000  1.0
# 20 B  744.0  0.5259  0.3835  0.0  0.000  0.500  1.0000  1.0     744.0  0.0806  0.0727  0.0  0.0000  0.0500  0.150  0.3500             744.0  0.2775  0.2671  0.0  0.0000  0.2267  0.4705  1.0
#    N  242.0  0.2624  0.3641  0.0  0.000  0.000  0.5000  1.0     242.0  0.0283  0.0391  0.0  0.0000  0.0000  0.050  0.2000             242.0  0.0836  0.1707  0.0  0.0000  0.0000  0.0903  1.0
# 50 B  744.0  0.7791  0.3356  0.0  0.600  1.000  1.0000  1.0     744.0  0.0423  0.0340  0.0  0.0200  0.0400  0.060  0.1800             744.0  0.2477  0.2309  0.0  0.0500  0.2000  0.3904  1.0
#    N  242.0  0.4732  0.3937  0.0  0.000  0.500  1.0000  1.0     242.0  0.0189  0.0193  0.0  0.0000  0.0200  0.020  0.0800             242.0  0.0892  0.1626  0.0  0.0000  0.0405  0.0906  1.0
# -  -  493.0  0.3670  0.3387  0.0  0.075  0.325  0.5833  1.0     493.0  0.0639  0.0767  0.0  0.0025  0.0512  0.085  0.4013             493.0  0.1772  0.2279  0.0  0.0062  0.1142  0.2552  1.0


# Random link masking - for each soident and semester single node
#     a_len        r_len     accuracy       recall
# count  6929.000000  6330.000000  6929.000000  6330.000000
# mean      5.001732     5.281991     0.934746     0.945108
# std       2.584540     2.506167     0.137306     0.162884
# min       1.000000     1.000000     0.000000     0.000000
# 25%       3.000000     3.000000     0.900000     1.000000
# 50%       5.000000     5.000000     1.000000     1.000000
# 75%       7.000000     7.000000     1.000000     1.000000
# max      17.000000    17.000000     1.000000     1.000000

# Random link masking - for each soident and each study type single node
#              a_len        r_len     accuracy       recall
# count  1215.000000  1178.000000  1215.000000  1178.000000
# mean      7.007407     7.146859     0.906007     0.932605
# std       3.089672     3.028271     0.132938     0.159386
# min       1.000000     1.000000     0.000000     0.000000
# 25%       5.000000     5.000000     0.833333     1.000000
# 50%       7.000000     7.000000     1.000000     1.000000
# 75%       9.000000     9.000000     1.000000     1.000000
# max      17.000000    17.000000     1.000000     1.000000


# Entire user masking - for each soident and each study type single node
#        a_len  r_len    accuracy      recall
# count  243.0  243.0  243.000000  243.000000
# mean   691.0  691.0    0.853591    0.960073
# std      0.0    0.0    0.055814    0.038947
# min    691.0  691.0    0.403763    0.666667
# 25%    691.0  691.0    0.832127    0.942582
# 50%    691.0  691.0    0.862518    0.963636
# 75%    691.0  691.0    0.888567    1.000000
# max    691.0  691.0    0.936324    1.000000
