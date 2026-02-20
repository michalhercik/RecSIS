import sys
from ast import literal_eval

import pandas as pd
import requests

sys.path.insert(0, "..")
from algo.embedder import EmbedderSyllabus
from data_repository import DataRepository
from progress import print_progress_bar
from scipy.spatial import distance

EMBEDDER_URL = "http://localhost:8003/embedding"
MEILI_LIMIT = 50
DATASET_FILE = "dataset.csv"
RND_STATE = 24
SAMPLE_FRAC = 0.1
PREFIX_MAGNIFICATION = 200

"""
Result
    - prefix has not effect on retrieved documents

PREFIX_MAGNIFICATION = 1
         embeddings  predictions
count  7.070000e+02        707.0
mean   4.051450e-16          0.0
std    6.210743e-15          0.0
min    0.000000e+00          0.0
25%    0.000000e+00          0.0
50%    0.000000e+00          0.0
75%    0.000000e+00          0.0
max    9.547918e-14          0.0

PREFIX_MAGNIFICATION = 50
         embeddings  predictions
count  7.070000e+02        707.0
mean   4.051450e-16          0.0
std    6.210743e-15          0.0
min    0.000000e+00          0.0
25%    0.000000e+00          0.0
50%    0.000000e+00          0.0
75%    0.000000e+00          0.0
max    9.547918e-14          0.0

PREFIX_MAGNIFICATION = 200
         embeddings  predictions
count  7.070000e+02        707.0
mean   4.051450e-16          0.0
std    6.210743e-15          0.0
min    0.000000e+00          0.0
25%    0.000000e+00          0.0
50%    0.000000e+00          0.0
75%    0.000000e+00          0.0
max    9.547918e-14          0.0
"""


def main():
    data = load_dataset()
    model = EmbedderSyllabus(DataRepository())
    model.fit()
    results = []
    data = data.sample(frac=SAMPLE_FRAC, random_state=RND_STATE)
    data = data.reset_index(drop=True)
    for i, sample in data.iterrows():
        if i % 5 == 0 or i == (data.shape[0] - 1):
            print_progress_bar(
                i, data.shape[0], prefix="Progress:", suffix="", length=25
            )
        prefix = PREFIX_MAGNIFICATION * create_prefix(sample)
        query1, predict1 = fetch(model, sample, "")
        query2, predict2 = fetch(model, sample, prefix)
        embeddings = create_embeddings([query1, query2])
        predictions = prediction_vectors(predict1, predict2)
        results.append(
            [
                distance.cosine(embeddings[0], embeddings[1]),
                distance.jaccard(predictions[0], predictions[1]),
            ]
        )

    df = pd.DataFrame(results, columns=["embeddings", "predictions"])
    print()
    print(df.describe())


def prediction_vectors(p1, p2):
    columns = list(set(p1 + p2))
    v1 = [1 if c in p1 else 0 for c in columns]
    v2 = [1 if c in p2 else 0 for c in columns]
    # print(p1)
    # print(p2)
    # print(pd.DataFrame([v1, v2], columns=columns))
    return [v1, v2]


def load_dataset():
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
    return df


def fetch(model, sample, prefix, limit=MEILI_LIMIT):
    query = model.build_query(sample["relevant"], prefix=prefix)
    filter_out = []
    filter_out.extend(sample["relevant"])
    filter_out.extend(sample["incompatible_with"])
    filter_out.extend(sample["interchange_for"])
    filter = model.build_filter(filter_out)
    predict = model.fetch_similar(query, filter, limit)
    return query, predict


def dot_product(a, b):
    """Compute dot product of two sequences (stops at shortest)."""
    return sum(x * y for x, y in zip(a, b))


def create_prefix(sample):
    prefix = "{0}-year student in a {1}’s program in {2}.".format(
        (
            "First"
            if sample["zroc"] == 1
            else "Second"
            if sample["zroc"] == 2
            else "Third"
        ),
        ("Bachelor" if sample["sdruh"] == "B" else "Master"),
        sample["sobor"],
    )
    return prefix


def create_embeddings(text):
    payload = {"model": "sbert", "text": text}
    resp = requests.post(
        EMBEDDER_URL,
        json=payload,
        headers={"Content-Type": "application/json"},
        timeout=10,
    )
    res = resp.json()
    # print(f"HTTP POST {EMBEDDER_URL} -> status: {resp.status_code}")
    embeddings = res.get("embedding")
    return embeddings


if __name__ == "__main__":
    main()
