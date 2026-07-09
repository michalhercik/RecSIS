import argparse
import sys

import numpy as np
import pandas as pd
from embedder import sbert_embed
from graphsage import dataset, eval, split

sys.path.insert(0, "..")
from data_repository import DataRepository


def main(args):
    VAL_RATIO = 0.2

    user, finished, povinn = dataset()
    train, val, test = split(finished, VAL_RATIO, 2024)

    val_results = (
        val.groupby("user_id")
        .agg({"course_id": list})
        .rename(columns={"course_id": "val_courses"})
        .reset_index()
    )

    test_results = (
        test.groupby("user_id")
        .agg({"course_id": list})
        .rename(columns={"course_id": "val_courses"})
        .reset_index()
    )

    results = (
        train.merge(povinn[["course_id", "embed"]], on="course_id")
        .groupby("user_id")
        .agg(
            {
                "user_id": "first",
                "course_id": set,
                "embed": lambda x: np.mean(x.values, axis=0),
            }
        )
        .rename(columns={"course_id": "train_courses"})
        .reset_index(drop=True)
    )

    if args.eval:
        print("eval")
        results = pd.merge(results, test_results, on="user_id")
    else:
        print("other")
        results = pd.merge(results, val_results, on="user_id")

    if args.mode == Mode.COURSE:
        results["pred"] = results.apply(lambda x: similar(povinn, x["embed"]), axis=1)
    elif args.mode == Mode.USER:
        results["pred"] = results.apply(
            lambda x: recommend_user(
                results[results["user_id"] != x["user_id"]].copy(), x["embed"]
            ),
            axis=1,
        )
    else:
        raise ValueError("Invalid mode")

    results["pred"] = results.apply(
        lambda x: [c for c in x["pred"] if c not in x["train_courses"]], axis=1
    )
    results["pred"] = results["pred"].apply(np.array)

    results["target"] = results.apply(
        lambda x: [1 if i in x["val_courses"] else 0 for i in x["pred"]], axis=1
    )
    results["target"] = results["target"].apply(np.array)

    results_description = eval(user, results)
    print(results_description)


def recommend_user(results, embed):
    def cos_sim(x1, x2):
        return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))

    results["sim"] = results.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
    results = results[["train_courses", "sim"]].explode("train_courses")
    results = results.sort_values("sim", ascending=False)
    results = results.drop_duplicates(subset="train_courses", keep="first")
    return results["train_courses"].to_list()


def similar(povinn, embed):
    def cos_sim(x1, x2):
        return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))

    povinn["sim"] = povinn.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
    povinn = povinn.sort_values("sim", ascending=False)
    return povinn["course_id"].to_list()


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
        "-m",
        "--mode",
        type=str,
        help="user/course mode",
    )
    return parser


class Mode:
    USER = "user"
    COURSE = "course"


if __name__ == "__main__":
    args = parser().parse_args()
    main(args)
