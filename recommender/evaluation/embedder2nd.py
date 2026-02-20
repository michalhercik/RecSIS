import argparse
import sys

import numpy as np
import pandas as pd
from embedder import sbert_embed
from graph import dataset, eval, split

sys.path.insert(0, "..")
from data_repository import DataRepository

# python .\embedder2nd.py -m course
#      Recall                                                      Precision                                        AP
#       count    mean     std  min     25%     50%     75%     max     count mean  std  min  25%  50%  75%  max  count    mean     std  min  25%     50%     75%     max
# 5  B  121.0  0.0133  0.0362  0.0  0.0000  0.0000  0.0000  0.2000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.0614  0.1792  0.0  0.0  0.0000  0.0000  1.0000
#    N   53.0  0.0311  0.0724  0.0  0.0000  0.0000  0.0000  0.3333      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1023  0.2512  0.0  0.0  0.0000  0.0000  1.0000
# 10 B  121.0  0.0325  0.0734  0.0  0.0000  0.0000  0.0526  0.5000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.0728  0.1647  0.0  0.0  0.0000  0.0000  1.0000
#    N   53.0  0.0470  0.0874  0.0  0.0000  0.0000  0.0667  0.3333      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1089  0.2477  0.0  0.0  0.0000  0.0500  1.0000
# 20 B  121.0  0.0780  0.1015  0.0  0.0000  0.0588  0.1250  0.6000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.0825  0.1301  0.0  0.0  0.0250  0.1000  0.5588
#    N   53.0  0.1075  0.1663  0.0  0.0000  0.0625  0.1429  0.6667      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1070  0.2121  0.0  0.0  0.0000  0.0990  1.0000
# 50 B  121.0  0.2018  0.1888  0.0  0.0833  0.1765  0.2857  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.0762  0.0899  0.0  0.0  0.0594  0.0968  0.5000
#    N   53.0  0.2429  0.2120  0.0  0.0625  0.2000  0.4000  0.7500      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.0909  0.1553  0.0  0.0  0.0485  0.0995  1.0000
# -  -   87.0  0.0943  0.1173  0.0  0.0182  0.0622  0.1341  0.5479      95.5  0.0  0.0  0.0  0.0  0.0  0.0  0.0   95.5  0.0878  0.1788  0.0  0.0  0.0166  0.0557  0.8824

# python .\embedder2nd.py -m user
#      Recall                                                      Precision                                        AP
#       count    mean     std  min     25%     50%     75%     max     count mean  std  min  25%  50%  75%  max  count    mean     std  min     25%     50%     75%     max
# 5  B  120.0  0.1308  0.2208  0.0  0.0000  0.0667  0.1667  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.2724  0.3327  0.0  0.0000  0.2000  0.5000  1.0000
#    N   52.0  0.0596  0.1357  0.0  0.0000  0.0000  0.0167  0.6667      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1194  0.2589  0.0  0.0000  0.0000  0.0000  1.0000
# 10 B  120.0  0.2499  0.2638  0.0  0.0000  0.1847  0.3636  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.2805  0.2764  0.0  0.0000  0.2460  0.4204  1.0000
#    N   52.0  0.1184  0.1941  0.0  0.0000  0.0000  0.1750  0.6667      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1347  0.2281  0.0  0.0000  0.0000  0.2083  1.0000
# 20 B  120.0  0.4286  0.2849  0.0  0.2308  0.4000  0.6038  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.2707  0.2251  0.0  0.0702  0.2472  0.4084  1.0000
#    N   52.0  0.2693  0.2699  0.0  0.0000  0.2071  0.4000  1.0000      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1515  0.2062  0.0  0.0000  0.0833  0.2019  1.0000
# 50 B  120.0  0.7303  0.2636  0.0  0.6429  0.8000  0.9333  1.0000     132.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  132.0  0.2452  0.1813  0.0  0.0779  0.2436  0.3510  0.8333
#    N   52.0  0.5262  0.2770  0.0  0.3636  0.5000  0.7330  1.0000      59.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0   59.0  0.1477  0.1635  0.0  0.0473  0.0854  0.1674  0.7500
# -  -   86.0  0.3141  0.2387  0.0  0.1547  0.2698  0.4240  0.9167      95.5  0.0  0.0  0.0  0.0  0.0  0.0  0.0   95.5  0.2028  0.2340  0.0  0.0244  0.1382  0.2822  0.9479


def main(args):
    VAL_RATIO = 0.2

    user, finished, povinn = dataset()
    train, val, test = split(finished, VAL_RATIO, 2024)

    pamela = DataRepository().pamela
    pamela = pamela[
        (pamela["jazyk"] == "ENG") & (pamela["typ"].isin(["A", "S"]))
    ].pivot_table(index="povinn", columns="typ", values="memo", aggfunc="first")
    povinn = povinn.merge(pamela, on="povinn")
    embed_src = povinn.apply(lambda x: f"{x['panazev']}: {x['A']}\n{x['S']}", axis=1)
    povinn["embed"] = list(sbert_embed(embed_src))

    val_results = (
        val.groupby("user_id")
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
    # TODO: inner join -> only users having both train and val courses included
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
        "-d",
        "--data",
        action="store_true",
        default=False,
        help="Load data from database",
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
