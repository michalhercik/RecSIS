import pandas as pd
from data import TrainData
from user import User
import psycopg2
import os

from recommender import EvalRecommender
from ranker.ranker import IdentityRanker
from categorizer import DepartmentCategorizer

ELSA = "Elsa"
GCN = "GCN"
LIGHT_GCN = "LightGCN"
CONTENT_KNN = "ContentKNN"
USER_KNN = "UserKNN"


ALGO = [ELSA, GCN, LIGHT_GCN, CONTENT_KNN, USER_KNN]
# ALGO = [GCN]


def main():
    user = User(
        "test-b1w",
        "NIPVS19B",
        2020,
        [
            {"year": 0, "unassigned": ["NPFL129"]},
            {"year": 1, "winter": ["NPRG021"], "summer": ["NDBI040"]},
        ],
    )

    data = TrainData(42)
    data.fit(cache=True)
    print(data.povinn.columns)

    ranker = IdentityRanker(data)
    ranker.fit()
    result = ranker.rank(user)

    categorizer = DepartmentCategorizer(data)
    cat_name, cat_values = categorizer.categorize(result)
    cat_values = [cat[:20] for cat in cat_values]
    for name, values in zip(cat_name, cat_values):
        print(f"{name}: {values}")
    exit()

    rec = EvalRecommender()
    rec.fit(ALGO)

    result = rec.recommend(user, ALGO, 10)
    print(result)

dsw = User(
    "test-dsw",
    "NIPVS19B",
    2020,
    [
        {
            "year": 0,
            "unassigned": [
                "NAIL062",
                "NDBI025",
                "NDMI002",
                "NDMI011",
                "NJAZ091",
                "NMAI054",
                "NMAI057",
                "NMAI058",
                "NMAI059",
                "NPRG013",
                "NPRG030",
                "NPRG031",
                "NPRG035",
                "NPRG041",
                "NPRG042",
                "NPRG043",
                "NPRG045",
                "NPRG054",
                "NPRG062",
                "NSWI004",
                "NSWI041",
                "NSWI098",
                "NSWI120",
                "NSWI141",
                "NSWI142",
                "NSWI154",
                "NSWI170",
                "NSWI177",
                "NTIN060",
                "NTIN061",
                "NTIN071",
                "NTVY014",
                "NTVY015",
                "NTVY016",
                "NTVY017",
            ],
        }
    ],
)
dai = User(
    "test-dai",
    "NIUI25B",
    2025,
    [
        {
            "year": 0,
            "unassigned": [
                "NAIL028",
                "NAIL062",
                "NAIL120",
                "NAIL121",
                "NDBI025",
                "NDMI002",
                "NDMI011",
                "NJAZ091",
                "NMAI054",
                "NMAI055",
                "NMAI057",
                "NMAI058",
                "NMAI059",
                "NPFL012",
                "NPFL124",
                "NPFL129",
                "NPGR036",
                "NPRG005",
                "NPRG013",
                "NPRG030",
                "NPRG031",
                "NPRG041",
                "NPRG045",
                "NPRG051",
                "NPRG062",
                "NSWI120",
                "NSWI141",
                "NSWI170",
                "NSWI177",
                "NTIN060",
                "NTIN061",
                "NTIN071",
                "NTVY014",
                "NTVY015",
                "NTVY016",
                "NTVY017",
            ],
        }
    ],
)
dgd = User(
    "test-dgd",
    "NIPGVAVH22B",
    2022,
    [
        {
            "year": 0,
            "unassigned": [
                "NAIL062",
                "NDBI025",
                "NDMI002",
                "NDMI011",
                "NJAZ091",
                "NMAI054",
                "NMAI055",
                "NMAI057",
                "NMAI058",
                "NMAI059",
                "NPGR002",
                "NPGR003",
                "NPGR025",
                "NPGR035",
                "NPGR037",
                "NPRG030",
                "NPRG031",
                "NPRG035",
                "NPRG041",
                "NPRG045",
                "NPRG062",
                "NSWI120",
                "NSWI141",
                "NSWI170",
                "NSWI177",
                "NTIN060",
                "NTIN061",
                "NTIN071",
                "NTVY014",
                "NTVY015",
                "NTVY016",
                "NTVY017",
            ],
        }
    ],
)

if __name__ == "__main__":
    main()
