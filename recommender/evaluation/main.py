import multiprocessing
import sys
import time
from ast import literal_eval
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime, timedelta

import pandas as pd
from scores import avg_prec_at_k, precision_at_k, recall_at_k
from sklearn.model_selection import train_test_split

sys.path.insert(0, "..")
from algo.embedder import Embedder
from data_repository import DataRepository

# Format: [name]_[dataset]-[y_true]
OUT_FILE_NAME_TAG = "searchable-interchange-d1-relevant"
OUT_FILE_NAME = datetime.now().strftime("%y%m%d-%H%M%S") + "-" + OUT_FILE_NAME_TAG
OUT_FILE = "eval/" + OUT_FILE_NAME + ".csv"
DESC_OUT_FILE = "desc/" + OUT_FILE_NAME + ".csv"
DATASET_FILE = "dataset.csv"
LIMIT = 50


def main():
    # y_true = ["1", "2", "3", "4", "9", "9", "9", "9"]
    # y_pred = ["5", "6", "7", "1", "2", "3"]
    # y_true = ["1", "2", "3"]
    # y_pred = ["1", "5", "2", "6", "3", "7", "4"]
    # k = 7
    # print(precision_at_k(y_true, y_pred, k))
    # print(recall_at_k(y_true, y_pred, k))
    # print(avg_prec_at_k(y_true, y_pred, k))
    # print(avg_prec_at_k(y_true, y_pred, k, normalize=True))
    # exit()
    df = pd.read_csv(
        DATASET_FILE,
        converters={
            "finished": literal_eval,
            "next_semester": lambda x: [] if len(x) == 0 else literal_eval(x),
            "relevant": lambda x: [] if len(x) == 0 else literal_eval(x),
        },
    )
    train, test = train_test_split(
        df["soident"].unique(), test_size=0.2, random_state=29980
    )
    do(exp(), df[df["soident"].isin(test)].reset_index(drop=True))
    # y_next_semester = sample["next_semester"]


class exp:
    def init(self):
        self.model = Embedder(DataRepository())
        self.model.fit()

    def get(self, X):
        # replace non searchable for searchable zamennosti
        query = self.model.build_query(X)
        # filter out zamennosti
        # filter out neslucitelnosti
        filter = self.model.build_filter(X)
        y_pred = self.model.fetch_similar(query, filter, LIMIT)
        return y_pred


def do(experiment, eval_data):
    start = time.time()
    experiment.init()

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
        y_pred = experiment.get(X)
        for k in [10, 20, 50]:
            results.append(
                [
                    id,
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
            "recommended",
        ],
    )
    out.to_csv(OUT_FILE, index=False)
    print()
    desc = (
        out.groupby("k")
        .describe(percentiles=[])[
            [
                ("rel_precision", "mean"),
                ("rel_precision", "std"),
                ("rel_recall", "mean"),
                ("rel_recall", "std"),
                ("rel_avg_prec", "mean"),
                ("rel_avg_prec", "std"),
                ("next_precision", "mean"),
                ("next_precision", "std"),
                ("next_recall", "mean"),
                ("next_recall", "std"),
                ("next_avg_prec", "mean"),
                ("next_avg_prec", "std"),
            ]
        ]
        .round(4)
    )
    desc.to_csv(DESC_OUT_FILE)
    print(desc)


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

# Equivalent courses
#
# select DISTINCT  p1.pnazev, p1.povinn, p1.pvyucovan, p2.pvyucovan, p2.povinn, p2.pnazev, preq1.reqtyp
# from preq preq1
# inner join preq preq2 on preq1.povinn = preq2.reqpovinn AND preq1.reqpovinn = preq2.povinn and preq1.reqtyp = preq2.reqtyp
# left join povinn p1 on preq1.povinn = p1.povinn
# left join povinn p2 on preq2.povinn = p2.povinn
# --left join povinn2searchable ps1 on p1.povinn = ps1.code
# --left join povinn2searchable ps2 on p2.povinn = ps2.code
# where preq1.reqtyp = 'Z'
# AND p2.pvyucovan = 'V'
# AND p2.pgarant not like '%STUD'
# --AND ps1.code is null
# --AND ps2.code is not null
