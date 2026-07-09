import argparse

import numpy as np
import pandas as pd
import torch
from elsa import ELSA
from graphsage import dataset, eval, split

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

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

    model = ELSA(
        n_items=povinn.shape[0], device=device, n_dims=args.factors, lr=args.lr
    )
    model.fit(
        train_im.values,
        batch_size=train.shape[0],
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

    if args.eval:
        results = test_results.merge(results, on="user_id")
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
    else:
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
        "-v",
        "--eval",
        action="store_true",
        default=False,
        help="Evaluate the model on the test set",
    )
    parser.add_argument(
        "-f", "--factors", type=int, default=16, help="Number of factors"
    )
    parser.add_argument(
        "-e",
        "--epochs",
        type=int,
        default=50,
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
