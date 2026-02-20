import argparse

import numpy as np
import pandas as pd
import torch
from graph import dataset, eval, split
from torch_geometric.nn.models import LightGCN
from torch_geometric.utils import negative_sampling

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# python .\lightgcn.py -e 200 -lr 1e-2
#      Recall                                                      Precision                                        AP
#       count    mean     std  min     25%     50%     75%     max     count mean  std  min  25%  50%  75%  max  count    mean     std  min     25%     50%     75%  max
# 5  B  132.0  0.1971  0.1399  0.0  0.0909  0.2000  0.3000  0.5000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.6082  0.3912  0.0  0.3062  0.6792  1.0000  1.0
#    N   59.0  0.1234  0.1238  0.0  0.0000  0.1111  0.1847  0.5000      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.4620  0.4221  0.0  0.0000  0.5000  1.0000  1.0
# 10 B  132.0  0.2893  0.1983  0.0  0.1405  0.2792  0.4500  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.5847  0.3607  0.0  0.3167  0.6467  0.9580  1.0
#    N   59.0  0.1593  0.1552  0.0  0.0000  0.1333  0.2222  0.6667      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.4491  0.3947  0.0  0.0000  0.3833  0.8500  1.0
# 20 B  132.0  0.4459  0.2379  0.0  0.2825  0.4545  0.6667  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.5027  0.3058  0.0  0.2552  0.5149  0.7512  1.0
#    N   59.0  0.2106  0.1712  0.0  0.0839  0.2000  0.3333  0.6667      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.4033  0.3415  0.0  0.0541  0.3833  0.6472  1.0
# 50 B  132.0  0.6032  0.2554  0.0  0.4286  0.6000  0.8000  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.4296  0.2649  0.0  0.2213  0.4190  0.6442  1.0
#    N   59.0  0.3249  0.2089  0.0  0.1742  0.3000  0.4833  0.8750      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.3384  0.2909  0.0  0.0654  0.2738  0.5000  1.0
# -  -   95.5  0.2942  0.1863  0.0  0.1501  0.2848  0.4300  0.7760      95.5  0.0  0.0  0.0  0.0  0.0  0.0  0.0   95.5  0.4722  0.3465  0.0  0.1524  0.4750  0.7938  1.0


def main(args):
    EPOCHS = args.epochs
    LR = 1e-2
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

    num_nodes = user.shape[0] + povinn.shape[0]
    model = LightGCN(num_nodes=num_nodes, embedding_dim=64, num_layers=3).to(device)

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
    parser.add_argument(
        "-lr",
        "--learning-rate",
        type=float,
        default=1e-2,
        help="Learning rate for the optimizer",
    )
    return parser


if __name__ == "__main__":
    args = parser().parse_args()
    main(args)
