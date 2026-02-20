import argparse

import numpy as np
import pandas as pd
import torch
from elsa import ELSA
from graph import dataset, eval, split

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# python .\elsa2nd.py -e 30 -lr 1e-2 -f 256
# Recall                                                         Precision                                        AP
#  count    mean     std     min     25%     50%     75%     max     count mean  std  min  25%  50%  75%  max  count    mean     std     min     25%     50%     75%     max
# 5  B  132.0  0.1671  0.1571  0.0000  0.0000  0.1429  0.2656  0.7500     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.4321  0.3467  0.0000  0.0000  0.4778  0.7000  1.0000
#    N   59.0  0.0533  0.0830  0.0000  0.0000  0.0000  0.0955  0.3333      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1626  0.2828  0.0000  0.0000  0.0000  0.2500  1.0000
# 10 B  132.0  0.3068  0.2233  0.0000  0.1667  0.3125  0.4000  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.4333  0.2894  0.0000  0.1991  0.4778  0.6429  1.0000
#    N   59.0  0.1464  0.1801  0.0000  0.0000  0.0909  0.2361  0.6667      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1747  0.2015  0.0000  0.0000  0.1250  0.2620  0.7667
# 20 B  132.0  0.4621  0.2564  0.0000  0.2875  0.4686  0.6250  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.3983  0.2389  0.0000  0.1998  0.4160  0.5638  1.0000
#    N   59.0  0.2722  0.2757  0.0000  0.0627  0.1667  0.4495  1.0000      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1826  0.1666  0.0000  0.0690  0.1429  0.2562  0.6429
# 50 B  132.0  0.6754  0.2281  0.0588  0.5417  0.6667  0.8333  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.3468  0.1949  0.0208  0.1734  0.3337  0.5138  0.8079
#    N   59.0  0.4842  0.3037  0.0000  0.2154  0.4545  0.7136  1.0000      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1664  0.1254  0.0000  0.0856  0.1347  0.2105  0.4983
# -  -   95.5  0.3209  0.2134  0.0074  0.1592  0.2878  0.4523  0.8438      95.5  0.0  0.0  0.0  0.0  0.0  0.0  0.0   95.5  0.2871  0.2308  0.0026  0.0909  0.2635  0.4249  0.8395


def main(args):
    VAL_RATIO = 0.2

    user, finished, povinn = dataset()
    train, val, test = split(finished, VAL_RATIO, 2024)

    def interaction_matrix(user, finished, povinn):
        im = pd.crosstab(finished["user_id"], finished["course_id"])
        im = im.reindex(
            index=user["user_id"], columns=povinn["course_id"], fill_value=0
        )
        return im

    train_im = interaction_matrix(user, train, povinn)
    # val_im = interaction_matrix(user, pd.concat([train, val]), povinn)
    # test_im = interaction_matrix(user, test, povinn)

    val_results = pd.merge(
        train.groupby("user_id")
        .agg({"course_id": set})
        .rename(columns={"course_id": "train_courses"}),
        val.groupby("user_id")
        .agg({"course_id": list})
        .rename(columns={"course_id": "val_courses"}),
        on="user_id",
    ).reset_index()

    model = ELSA(
        n_items=povinn.shape[0], device=device, n_dims=args.factors, lr=args.lr
    )
    model.fit(
        train_im.values,
        batch_size=train.shape[0],  # TODO: Batch size of train data size
        epochs=args.epochs,
        shuffle=False,
    )

    predictions = model.predict(train_im.values, batch_size=test.shape[0])

    results = pd.DataFrame(
        {
            "user_id": train_im.index.tolist(),
            "pred": map(
                lambda indices: np.array(train_im.columns[indices].to_list()),
                torch.topk(predictions, k=predictions.shape[1], sorted=True).indices,
            ),
        }
    )
    results = val_results.merge(results, on="user_id")
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
        "-f", "--factors", type=int, default=256, help="Number of factors"
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
    return parser


if __name__ == "__main__":
    args = parser().parse_args()
    main(args)
