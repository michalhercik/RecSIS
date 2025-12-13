import sys
import time
from ast import literal_eval
from datetime import datetime, timedelta

import pandas as pd
from scores import avg_prec_at_k, precision_at_k, recall_at_k
from sklearn.model_selection import train_test_split

sys.path.insert(0, "..")
from data_repository import DataRepository
from elsaexp import ElsaExp
from embedder import (
    EmbedderAnnotationExp,
    EmbedderExp1,
    EmbedderExp2,
    EmbedderExp3,
    EmbedderSyllabusExp,
)

DATASET_FILE = "dataset.csv"
OUT_PATH = "eval"
RND_STATE = 29980
TEST_SIZE = 0.1
RETRIEVE_LIMIT = 50


def main():
    df = pd.read_csv(
        DATASET_FILE,
        converters={
            "finished": literal_eval,
            "next_semester": lambda x: [] if len(x) == 0 else literal_eval(x),
            "relevant": lambda x: [] if len(x) == 0 else literal_eval(x),
            "interchange_for": lambda x: [] if len(x) == 0 else literal_eval(x),
            "incompatible_with": lambda x: [] if len(x) == 0 else literal_eval(x),
        },
    )
    train, test = train_test_split(
        df["soident"].unique(), test_size=TEST_SIZE, random_state=RND_STATE
    )
    train, validate = train_test_split(
        train, test_size=TEST_SIZE, random_state=RND_STATE
    )
    # print(train.shape)
    # print(validate.shape)
    # print(test.shape)
    data_repository = DataRepository()
    train_data = df[df["soident"].isin(train)].reset_index(drop=True)
    validate_data = df[df["soident"].isin(validate)].reset_index(drop=True)
    test_data = df[df["soident"].isin(test)].reset_index(drop=True)
    print(train_data.shape)
    print(validate_data.shape)
    print(test_data.shape)
    exit()
    embedder_experiments = [
        # EmbedderExp1(data_repository),
        # EmbedderExp2(data_repository),
        # EmbedderExp3(data_repository),
        # EmbedderAnnotationExp(data_repository),
        # EmbedderSyllabusExp(data_repository),
    ]
    elsa_experiments = [ElsaExp(data_repository, train_data, validate_data)]
    for exp in elsa_experiments:
        print()
        print(50 * "=")
        print(exp.__class__.__name__)
        print(50 * "-")
        elsa_do(exp, train_data, safe=False)
        elsa_do(exp, validate_data, safe=False)
        elsa_do(exp, test_data)
    for exp in embedder_experiments:
        print()
        print(50 * "=")
        print(exp.__class__.__name__)
        print(50 * "-")
        do(exp, test_data)


def elsa_do(exp, data, safe=True):
    y_pred = exp.get(data, 50)
    results = []
    for i, sample in data.iterrows():
        pred = y_pred[i]
        for y_true_label in ["next_semester", "relevant"]:
            y_true = sample[y_true_label]
            for k in [3, 10, 20, 50]:
                results.append(
                    [
                        id,
                        sample["sdruh"],
                        sample["zroc"],
                        sample["zsem"],
                        y_true_label,
                        k,
                        len(y_true) / k,
                        precision_at_k(y_true, pred, k),
                        recall_at_k(y_true, pred, k),
                        avg_prec_at_k(y_true, pred, k),
                        pred,
                    ]
                )
    out = pd.DataFrame(
        data=results,
        columns=[
            "id",
            "sdruh",
            "zroc",
            "zsem",
            "target",
            "k",
            "true/k",
            "precision",
            "recall",
            "avg_prec",
            "predicted",
        ],
    )
    out.to_csv(out_file_path(exp.name), index=False)
    desc = (
        out.groupby(["target", "k"])
        .describe()[
            [
                ("precision", "mean"),
                ("precision", "std"),
                ("precision", "25%"),
                ("precision", "50%"),
                ("precision", "75%"),
                ("recall", "mean"),
                ("recall", "std"),
                ("recall", "25%"),
                ("recall", "50%"),
                ("recall", "75%"),
                ("avg_prec", "mean"),
                ("avg_prec", "std"),
                ("avg_prec", "25%"),
                ("avg_prec", "50%"),
                ("avg_prec", "75%"),
            ]
        ]
        .round(2)
    )
    print(desc)


def do(experiment, eval_data):
    start = time.time()

    results = []
    for i, sample in eval_data.iterrows():
        if i % 5 == 0 or i == eval_data.shape[0]:
            print_progress_bar(
                i, eval_data.shape[0], prefix="Progress:", suffix="", length=25
            )
        X = sample["finished"]
        relevant_true = sample["relevant"]
        next_true = sample["next_semester"]
        id = sample["soident"]
        # prefix = "{0}-year student in a {1}’s program in {2}.".format(
        #     (
        #         "First"
        #         if sample["zroc"] == 1
        #         else "Second"
        #         if sample["zroc"] == 2
        #         else "Third"
        #     ),
        #     ("Bachelor" if sample["sdruh"] == "B" else "Master"),
        #     sample["sobor"],
        # )
        try:
            y_pred = experiment.get(
                X,
                sample["incompatible_with"],
                sample["interchange_for"],
                RETRIEVE_LIMIT,
            )
        except Exception as e:
            print()
            print(e)
            print(i)
            exit()
        for k in [3, 10, 20, 50]:
            results.append(
                [
                    id,
                    sample["sdruh"],
                    sample["zroc"],
                    sample["zsem"],
                    k,
                    len(relevant_true) / k,
                    precision_at_k(relevant_true, y_pred, k),
                    recall_at_k(relevant_true, y_pred, k),
                    avg_prec_at_k(relevant_true, y_pred, k),
                    len(next_true) / k,
                    precision_at_k(next_true, y_pred, k),
                    recall_at_k(next_true, y_pred, k),
                    avg_prec_at_k(next_true, y_pred, k),
                    X,
                    relevant_true,
                    next_true,
                    y_pred,
                ]
            )

    exec_time = time.time() - start
    print()
    print("Execution time:", timedelta(seconds=exec_time))
    out = pd.DataFrame(
        data=results,
        columns=[
            "id",
            "sdruh",
            "zroc",
            "zsem",
            "k",
            "rel_true/k",
            "rel_precision",
            "rel_recall",
            "rel_avg_prec",
            "next_true/k",
            "next_precision",
            "next_recall",
            "next_avg_prec",
            "finished",
            "relevant_true",
            "next_true",
            "predicted",
        ],
    )
    out.to_csv(out_file_path(experiment.name), index=False)
    print()
    desc = (
        out.groupby("k")
        .describe(percentiles=[0.75])[
            [
                ("rel_precision", "mean"),
                ("rel_precision", "std"),
                ("rel_precision", "75%"),
                ("rel_recall", "mean"),
                ("rel_recall", "std"),
                ("rel_recall", "75%"),
                ("rel_avg_prec", "mean"),
                ("rel_avg_prec", "std"),
                ("rel_avg_prec", "75%"),
                ("next_precision", "mean"),
                ("next_precision", "std"),
                ("next_precision", "75%"),
                ("next_recall", "mean"),
                ("next_recall", "std"),
                ("next_recall", "75%"),
                ("next_avg_prec", "mean"),
                ("next_avg_prec", "std"),
                ("next_avg_prec", "75%"),
            ]
        ]
        .round(2)
    )
    print(desc)


def out_file_path(name):
    datetime_tag = datetime.now().strftime("%y%m%d-%H%M%S")
    return f"{OUT_PATH}/{name}-{datetime_tag}.csv"


# Print iterations progress
def print_progress_bar(
    iteration,
    total,
    prefix="",
    suffix="",
    decimals=1,
    length=100,
    fill="█",
    printEnd="\r",
):
    """
    Call in a loop to create terminal progress bar
    @params:
        iteration   - Required  : current iteration (Int)
        total       - Required  : total iterations (Int)
        prefix      - Optional  : prefix string (Str)
        suffix      - Optional  : suffix string (Str)
        decimals    - Optional  : positive number of decimals in percent complete (Int)
        length      - Optional  : character length of bar (Int)
        fill        - Optional  : bar fill character (Str)
        printEnd    - Optional  : end character (e.g. "\r", "\r\n") (Str)
    """
    percent = ("{0:." + str(decimals) + "f}").format(100 * (iteration / float(total)))
    filledLength = int(length * iteration // total)
    bar = fill * filledLength + "-" * (length - filledLength)
    print(f"\r{prefix} |{bar}| {iteration}/{total} {suffix}", end=printEnd)
    # Print New Line on Complete
    if iteration == total:
        print()


if __name__ == "__main__":
    main()
