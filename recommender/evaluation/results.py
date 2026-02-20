import sys
from ast import literal_eval
from sqlite3.dbapi2 import converters

import pandas as pd
from scores import avg_prec_at_k, precision_at_k, recall_at_k

sys.path.insert(0, "..")
from data_repository import DataRepository

# ELSA_EVAL = "ElsaRelevant-251215-210051"
# ELSA_EVAL = "ElsaFinishedBachelorsExp-251221-155626"
ELSA_EVAL = "ElsaFinishedMastersExp-251221-155639"

K = 10


def main():
    elsa_eval_without_mandatory()


def elsa_eval():
    """
                        precision                                       recall                                     avg_prec
                        count  mean   std  min  25%  50%   75%  max   mean   std   min   25%   50%   75%   max     mean   std   min   25%   50%   75%   max
    sdruh zroc zsem
    B     1    1         75.0  0.98  0.05  0.8  1.0  1.0  1.00  1.0   0.24  0.06  0.11  0.20  0.25  0.28  0.38     0.98  0.05  0.77  1.00  1.00  1.00  1.00
               2         75.0  0.93  0.10  0.6  0.9  1.0  1.00  1.0   0.29  0.08  0.11  0.23  0.31  0.34  0.47     0.95  0.08  0.66  0.93  1.00  1.00  1.00
          2    1         75.0  0.78  0.16  0.3  0.7  0.8  0.90  1.0   0.32  0.10  0.13  0.25  0.33  0.39  0.64     0.87  0.13  0.38  0.82  0.90  0.98  1.00
               2         75.0  0.54  0.25  0.0  0.4  0.5  0.70  1.0   0.31  0.14  0.00  0.24  0.30  0.39  0.58     0.65  0.27  0.00  0.48  0.69  0.89  1.00
          3    1         75.0  0.37  0.31  0.0  0.1  0.3  0.60  1.0   0.32  0.24  0.00  0.17  0.31  0.46  1.00     0.48  0.36  0.00  0.11  0.50  0.83  1.00
               2         43.0  0.46  0.38  0.0  0.0  0.5  0.80  1.0   0.27  0.23  0.00  0.00  0.29  0.43  1.00     0.53  0.42  0.00  0.00  0.63  0.93  1.00
          4    1         37.0  0.46  0.40  0.0  0.0  0.5  0.90  1.0   0.25  0.20  0.00  0.00  0.29  0.42  0.67     0.53  0.43  0.00  0.00  0.64  0.98  1.00
               2         23.0  0.71  0.29  0.0  0.6  0.8  0.90  1.0   0.34  0.15  0.00  0.27  0.37  0.43  0.64     0.80  0.29  0.00  0.73  0.91  1.00  1.00
          5    1         21.0  0.78  0.19  0.4  0.6  0.8  0.90  1.0   0.37  0.10  0.22  0.29  0.39  0.43  0.64     0.87  0.14  0.60  0.74  0.95  1.00  1.00
               2         21.0  0.78  0.19  0.4  0.6  0.8  0.90  1.0   0.37  0.10  0.22  0.29  0.39  0.43  0.64     0.87  0.14  0.60  0.74  0.95  1.00  1.00
          6    1         21.0  0.78  0.19  0.4  0.6  0.8  0.90  1.0   0.37  0.10  0.22  0.29  0.39  0.43  0.64     0.87  0.14  0.60  0.74  0.95  1.00  1.00
               2         21.0  0.78  0.19  0.4  0.6  0.8  0.90  1.0   0.37  0.10  0.22  0.29  0.39  0.43  0.64     0.87  0.14  0.60  0.74  0.95  1.00  1.00
    N     1    1         45.0  0.46  0.21  0.0  0.3  0.5  0.60  0.9   0.39  0.18  0.00  0.29  0.38  0.46  0.78     0.62  0.26  0.00  0.47  0.67  0.83  1.00
               2         43.0  0.17  0.15  0.0  0.0  0.2  0.30  0.5   0.28  0.24  0.00  0.00  0.25  0.41  0.83     0.31  0.29  0.00  0.00  0.25  0.58  1.00
          2    1         35.0  0.07  0.13  0.0  0.0  0.0  0.10  0.5   0.18  0.35  0.00  0.00  0.00  0.11  1.00     0.08  0.16  0.00  0.00  0.00  0.12  0.57
               2          7.0  0.09  0.19  0.0  0.0  0.0  0.05  0.5   0.18  0.37  0.00  0.00  0.00  0.12  1.00     0.11  0.22  0.00  0.00  0.00  0.10  0.57
          3    1          5.0  0.08  0.18  0.0  0.0  0.0  0.00  0.4   0.20  0.45  0.00  0.00  0.00  0.00  1.00     0.09  0.20  0.00  0.00  0.00  0.00  0.46
    """
    df = pd.read_csv(f"eval/{ELSA_EVAL}.csv")
    # df["true_count"] = df["true/k"] * df["k"]
    df = df.drop(["soident", "true/k", "predicted"], axis=1)
    print(df.columns)
    df = (
        df[(df["target"] == "relevant") & (df["k"] == K)]
        .drop(["k", "target"], axis=1)
        .groupby(["sdruh", "zroc", "zsem"])
    )
    print(
        df.describe()
        .round(2)
        .drop(
            [("recall", "count"), ("avg_prec", "count"), ("true_count", "count")],
            axis=1,
            errors="ignore",
        )
    )


def elsa_eval_without_mandatory():
    """
    Without filtering out data with no relevant courses
                        precision                                       recall                                    avg_prec
                        count  mean   std  min  25%  50%   75%  max   mean   std  min   25%   50%   75%   max     mean   std   min   25%   50%   75%   max
    sdruh zroc zsem
    B     1    1         75.0  0.52  0.33  0.0  0.2  0.5  0.90  1.0   0.37  0.19  0.0  0.23  0.33  0.50  1.00     0.70  0.30  0.00  0.50  0.83  1.00  1.00
               2         75.0  0.44  0.35  0.0  0.2  0.3  0.80  1.0   0.34  0.24  0.0  0.19  0.33  0.50  1.00     0.62  0.38  0.00  0.28  0.75  0.99  1.00
          2    1         75.0  0.38  0.35  0.0  0.1  0.3  0.70  1.0   0.31  0.26  0.0  0.12  0.27  0.44  1.00     0.51  0.39  0.00  0.14  0.50  0.91  1.00
               2         75.0  0.31  0.34  0.0  0.0  0.2  0.60  1.0   0.26  0.29  0.0  0.00  0.23  0.40  1.00     0.42  0.40  0.00  0.00  0.42  0.85  1.00
          3    1         75.0  0.29  0.34  0.0  0.0  0.1  0.60  1.0   0.24  0.27  0.0  0.00  0.21  0.39  1.00     0.37  0.41  0.00  0.00  0.17  0.82  1.00
               2         43.0  0.45  0.37  0.0  0.1  0.4  0.80  1.0   0.33  0.26  0.0  0.19  0.30  0.44  1.00     0.56  0.41  0.00  0.12  0.74  0.94  1.00
          4    1         37.0  0.44  0.39  0.0  0.0  0.5  0.80  1.0   0.26  0.24  0.0  0.00  0.27  0.43  1.00     0.52  0.43  0.00  0.00  0.64  0.93  1.00
               2         23.0  0.70  0.28  0.0  0.6  0.7  0.90  1.0   0.35  0.15  0.0  0.28  0.37  0.44  0.64     0.79  0.29  0.00  0.76  0.90  0.99  1.00
          5    1         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
               2         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
          6    1         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
               2         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
    N     1    1         45.0  0.14  0.20  0.0  0.0  0.1  0.20  0.7   0.26  0.30  0.0  0.00  0.20  0.50  1.00     0.28  0.34  0.00  0.00  0.12  0.42  1.00
               2         43.0  0.06  0.12  0.0  0.0  0.0  0.10  0.6   0.11  0.20  0.0  0.00  0.00  0.19  0.67     0.13  0.24  0.00  0.00  0.00  0.18  0.79
          2    1         35.0  0.02  0.05  0.0  0.0  0.0  0.00  0.2   0.08  0.24  0.0  0.00  0.00  0.00  1.00     0.03  0.09  0.00  0.00  0.00  0.00  0.42
               2          7.0  0.03  0.05  0.0  0.0  0.0  0.05  0.1   0.06  0.10  0.0  0.00  0.00  0.08  0.25     0.03  0.06  0.00  0.00  0.00  0.05  0.12
          3    1          5.0  0.00  0.00  0.0  0.0  0.0  0.00  0.0   0.00  0.00  0.0  0.00  0.00  0.00  0.00     0.00  0.00  0.00  0.00  0.00  0.00  0.00

    Filtering out data with no relevant courses
                        precision                                       recall                                    avg_prec
                        count  mean   std  min  25%  50%   75%  max   mean   std  min   25%   50%   75%   max     mean   std   min   25%   50%   75%   max
    sdruh zroc zsem
    B     1    1         75.0  0.52  0.33  0.0  0.2  0.5  0.90  1.0   0.37  0.19  0.0  0.23  0.33  0.50  1.00     0.70  0.30  0.00  0.50  0.83  1.00  1.00
               2         74.0  0.45  0.35  0.0  0.2  0.3  0.80  1.0   0.35  0.23  0.0  0.20  0.33  0.50  1.00     0.63  0.37  0.00  0.29  0.76  1.00  1.00
          2    1         71.0  0.40  0.35  0.0  0.1  0.3  0.70  1.0   0.33  0.26  0.0  0.15  0.28  0.48  1.00     0.53  0.38  0.00  0.18  0.55  0.92  1.00
               2         70.0  0.34  0.34  0.0  0.0  0.2  0.68  1.0   0.28  0.29  0.0  0.00  0.24  0.41  1.00     0.45  0.39  0.00  0.00  0.43  0.87  1.00
          3    1         62.0  0.35  0.35  0.0  0.0  0.2  0.70  1.0   0.29  0.27  0.0  0.00  0.26  0.47  1.00     0.44  0.41  0.00  0.00  0.41  0.88  1.00
               2         43.0  0.45  0.37  0.0  0.1  0.4  0.80  1.0   0.33  0.26  0.0  0.19  0.30  0.44  1.00     0.56  0.41  0.00  0.12  0.74  0.94  1.00
          4    1         34.0  0.48  0.38  0.0  0.0  0.6  0.80  1.0   0.28  0.24  0.0  0.00  0.28  0.43  1.00     0.56  0.42  0.00  0.00  0.70  0.96  1.00
               2         23.0  0.70  0.28  0.0  0.6  0.7  0.90  1.0   0.35  0.15  0.0  0.28  0.37  0.44  0.64     0.79  0.29  0.00  0.76  0.90  0.99  1.00
          5    1         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
               2         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
          6    1         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
               2         21.0  0.76  0.19  0.4  0.6  0.8  0.90  1.0   0.38  0.11  0.2  0.32  0.40  0.44  0.64     0.87  0.15  0.44  0.78  0.92  1.00  1.00
    N     1    1         41.0  0.16  0.20  0.0  0.0  0.1  0.20  0.7   0.29  0.30  0.0  0.00  0.25  0.50  1.00     0.30  0.35  0.00  0.00  0.17  0.46  1.00
               2         32.0  0.08  0.13  0.0  0.0  0.0  0.10  0.6   0.15  0.22  0.0  0.00  0.00  0.33  0.67     0.18  0.27  0.00  0.00  0.00  0.27  0.79
          2    1         19.0  0.04  0.07  0.0  0.0  0.0  0.05  0.2   0.14  0.32  0.0  0.00  0.00  0.05  1.00     0.06  0.12  0.00  0.00  0.00  0.06  0.42
               2          5.0  0.04  0.05  0.0  0.0  0.0  0.10  0.1   0.08  0.12  0.0  0.00  0.00  0.17  0.25     0.04  0.06  0.00  0.00  0.00  0.10  0.12
          3    1          3.0  0.00  0.00  0.0  0.0  0.0  0.00  0.0   0.00  0.00  0.0  0.00  0.00  0.00  0.00     0.00  0.00  0.00  0.00  0.00  0.00  0.00
    """
    df = pd.read_csv(
        f"eval/{ELSA_EVAL}.csv",
        converters={"predicted": lambda x: [] if len(x) == 0 else literal_eval(x)},
    )
    dataset = pd.read_csv(
        "dataset.csv",
        usecols=["soident", "sdruh", "zroc", "zsem", "relevant"],
        converters={
            "relevant": lambda x: [] if len(x) == 0 else literal_eval(x),
        },
    )
    data_repo = DataRepository()
    sp = data_repo.stud_plan[["plan_code", "plan_year", "code"]]
    sp = sp.groupby(["plan_code", "plan_year"]).agg(list).reset_index()
    sp["plan_year"] = sp["plan_year"].astype(str)
    sp = pd.merge(
        data_repo.studium[["soident", "sdruh", "splan", "srokp"]],
        sp,
        left_on=["splan", "srokp"],
        right_on=["plan_code", "plan_year"],
    )
    sp = sp.drop(columns=["plan_code", "plan_year", "srokp", "splan"])

    df = df[(df["target"] == "relevant") & (df["k"] == K)].drop(["k", "target"], axis=1)
    df = pd.merge(df, sp, on=["soident", "sdruh"], how="left")
    df = df.rename(columns={"code": "splan"})
    df = df[["soident", "sdruh", "zroc", "zsem", "predicted", "splan"]]
    df = df.drop_duplicates(subset=["soident", "sdruh", "zroc", "zsem"])

    df = pd.merge(df, dataset, on=["soident", "sdruh", "zroc", "zsem"])

    df = remove_splan(df, ["predicted", "relevant"])

    # If no relevant courses availabe -> there are no good predictions
    df = df[df["relevant"].apply(len) > 0]

    df["precision"] = df.apply(
        lambda x: precision_at_k(x["relevant"], x["predicted"], K), axis=1
    )
    df["recall"] = df.apply(
        lambda x: recall_at_k(x["relevant"], x["predicted"], K), axis=1
    )
    df["avg_prec"] = df.apply(
        lambda x: avg_prec_at_k(x["relevant"], x["predicted"], K), axis=1
    )
    df = df.drop(["soident", "predicted", "splan", "relevant"], axis=1).groupby(
        ["sdruh", "zroc", "zsem"]
    )
    print(
        df.describe()
        .round(2)
        .drop(
            [("recall", "count"), ("avg_prec", "count"), ("true_count", "count")],
            axis=1,
            errors="ignore",
        )
    )
    pass


def remove_splan(df, labels):
    df = df.copy()
    df["splan"] = df["splan"].apply(lambda d: d if isinstance(d, list) else [])

    for label in labels:
        df[label] = df[[label, "splan"]].apply(
            lambda x: [e for e in x[label] if e not in x["splan"]],
            axis=1,
        )
    return df


def embedder_eval():
    df = pd.read_csv("eval/EmbedderSyllabus-251205-192543.csv")
    # df = pd.read_csv("eval/EmbedderSyllabus-251206-144529.csv")
    # print(
    #     df[(df["rel_avg_prec"] == 1) & (df["k"] == 10) & (df["sdruh"] == "B")][
    #         ["sdruh", "zroc", "zsem", "rel_precision", "rel_recall", "rel_avg_prec"]
    #     ]
    #     .describe()
    #     .round(2)
    # )
    print(
        df[df["sdruh"] == "B"]
        .groupby("k")[
            [
                "rel_precision",
                "rel_recall",
                "rel_avg_prec",
                "next_precision",
                "next_recall",
                "next_avg_prec",
            ]
        ]
        .describe(percentiles=[0.8])[
            [
                ("rel_precision", "mean"),
                ("rel_precision", "std"),
                ("rel_precision", "80%"),
                ("rel_recall", "mean"),
                ("rel_recall", "std"),
                ("rel_recall", "80%"),
                ("rel_avg_prec", "mean"),
                ("rel_avg_prec", "std"),
                ("rel_avg_prec", "80%"),
                ("next_precision", "mean"),
                ("next_precision", "std"),
                ("next_precision", "80%"),
                ("next_recall", "mean"),
                ("next_recall", "std"),
                ("next_recall", "80%"),
                ("next_avg_prec", "mean"),
                ("next_avg_prec", "std"),
                ("next_avg_prec", "80%"),
            ]
        ]
        .round(2)
    )
    print(
        df[df["sdruh"] == "N"]
        .groupby("k")[
            [
                "rel_precision",
                "rel_recall",
                "rel_avg_prec",
                "next_precision",
                "next_recall",
                "next_avg_prec",
            ]
        ]
        .describe(percentiles=[0.8])[
            [
                ("rel_precision", "mean"),
                ("rel_precision", "std"),
                ("rel_precision", "80%"),
                ("rel_recall", "mean"),
                ("rel_recall", "std"),
                ("rel_recall", "80%"),
                ("rel_avg_prec", "mean"),
                ("rel_avg_prec", "std"),
                ("rel_avg_prec", "80%"),
                ("next_precision", "mean"),
                ("next_precision", "std"),
                ("next_precision", "80%"),
                ("next_recall", "mean"),
                ("next_recall", "std"),
                ("next_recall", "80%"),
                ("next_avg_prec", "mean"),
                ("next_avg_prec", "std"),
                ("next_avg_prec", "80%"),
            ]
        ]
        .round(2)
    )


def eval(df):
    results = []
    for i, sample in df.iterrows():
        relevant_true = sample["relevant_true"]
        next_true = sample["next_true"]
        y_pred = sample["predicted"]
        k = sample["k"]
        for k in [3, 10, 20, 50]:
            results.append(
                [
                    sample["id"],
                    sample["sdruh"],
                    sample["zroc"],
                    sample["zsem"],
                    k,
                    len(relevant_true) / k,
                    # precision_at_k(relevant_true, y_pred, k),
                    # recall_at_k(relevant_true, y_pred, k),
                    # avg_prec_at_k(relevant_true, y_pred, k),
                    # len(next_true) / k,
                    # precision_at_k(next_true, y_pred, k),
                    # recall_at_k(next_true, y_pred, k),
                    # avg_prec_at_k(next_true, y_pred, k),
                ]
            )


if __name__ == "__main__":
    main()
