import pandas as pd


def main():
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
