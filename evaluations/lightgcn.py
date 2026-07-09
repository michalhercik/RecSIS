import argparse

import numpy as np
import pandas as pd
import torch
from graphsage import dataset, eval, split
from torch_geometric.nn.models import LightGCN
from torch_geometric.utils import negative_sampling

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

def main(args):
    EPOCHS = args.epochs
    LR = args.lr
    VAL_RATIO = 0.2

    user, finished, povinn = dataset()
    povinn["course_id"] = povinn["course_id"] + user.shape[0]
    finished["course_id"] = finished["course_id"] + user.shape[0]
    train, val, test = split(finished, VAL_RATIO, 2024)

    val_results = pd.merge(
        train.groupby("user_id")
        .agg({"course_id": set})
        .rename(columns={"course_id": "train_courses"}),
        val.groupby("user_id")
        .agg({"course_id": list})
        .rename(columns={"course_id": "val_courses"}),
        on="user_id",
    ).reset_index()

    test_results = pd.merge(
        train.groupby("user_id")
        .agg({"course_id": set})
        .rename(columns={"course_id": "train_courses"}),
        test.groupby("user_id")
        .agg({"course_id": list})
        .rename(columns={"course_id": "val_courses"}),
        on="user_id",
    ).reset_index()

    num_nodes = user.shape[0] + povinn.shape[0]
    model = LightGCN(num_nodes=num_nodes, embedding_dim=16, num_layers=2).to(device)

    edge_index_homo = torch.stack(
        [
            torch.tensor(train["user_id"].values),
            torch.tensor(train["course_id"].values),
        ],
        dim=0,
    )
    edge_index_homo = torch.cat([edge_index_homo, edge_index_homo.flip(0)], dim=1)

    optimizer = torch.optim.AdamW(model.parameters(), lr=LR)

    def train_step():
        model.train()
        optimizer.zero_grad()

        neg_edge_index = negative_sampling(
            edge_index=edge_index_homo,
            num_nodes=num_nodes,
            num_neg_samples=edge_index_homo.size(1) // 2,
        )

        pos_u, pos_i = edge_index_homo[:, : train.shape[0]]
        _, neg_i = neg_edge_index

        emb = model.get_embedding(edge_index_homo)

        u_emb = emb[pos_u]
        pos_emb = emb[pos_i]
        neg_emb = emb[neg_i]

        pos_scores = (u_emb * pos_emb).sum(dim=1)
        neg_scores = (u_emb * neg_emb).sum(dim=1)

        loss = model.recommendation_loss(
            pos_scores,
            neg_scores,
            node_id=torch.cat([pos_u, pos_i, neg_i]),
            lambda_reg=1e-4,
        )
        loss.backward()
        optimizer.step()

        return float(loss.detach())

    for epoch in range(1, EPOCHS + 1):
        loss = train_step()
        print(f"Epoch {epoch:03d} | Loss: {loss:.4f}")

    model.eval()

    if args.eval:
        top_items = model.recommend(
            edge_index=edge_index_homo,
            src_index=torch.tensor(test_results["user_id"].values),
            k=povinn.shape[0],
        )
        results = pd.merge(
            test_results,
            pd.DataFrame({"user_id": test_results["user_id"], "pred": top_items.tolist()}),
            on=["user_id"],
        )
    else:
        top_items = model.recommend(
            edge_index=edge_index_homo,
            src_index=torch.tensor(val_results["user_id"].values),
            k=povinn.shape[0],
        )
        results = pd.merge(
            val_results,
            pd.DataFrame({"user_id": val_results["user_id"], "pred": top_items.tolist()}),
            on=["user_id"],
        )

    results["pred"] = results.apply(
        lambda x: [i for i in x["pred"] if i not in x["train_courses"]], axis=1
    )
    results["pred"] = results["pred"].apply(np.array)
    results["target"] = results.apply(
        lambda x: [1 if i in x["val_courses"] else 0 for i in x["pred"]], axis=1
    )
    results["target"] = results["target"].apply(np.array)

    results_description = eval(user, results)
    print(results_description)


def parser():
    parser = argparse.ArgumentParser(
        prog="RECSIS Graph Recommender",
        description="Train and evaluate the model",
    )
    parser.add_argument(
        "-e",
        "--epochs",
        type=int,
        default=10,
        help="Number of epochs to train the model",
    )
    parser.add_argument(
        "-v",
        "--eval",
        action="store_true",
        default=False,
        help="Evaluate the model on the test set",
    )
    parser.add_argument(
        "-lr",
        "--learning-rate",
        type=float,
        default=1e-2,
        dest="lr",
        help="Learning rate for the optimizer",
    )
    return parser


if __name__ == "__main__":
    args = parser().parse_args()
    main(args)
