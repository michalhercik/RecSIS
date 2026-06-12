import pandas as pd
from data import TrainData
from user import User

from recommender import EvalRecommender

ELSA = "Elsa"
GCN = "GCN"
LIGHT_GCN = "LightGCN"
CONTENT_KNN = "ContentKNN"
USER_KNN = "UserKNN"


ALGO = [CONTENT_KNN]


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
    data.fit()

    rec = EvalRecommender()
    rec.fit(ALGO)
    result = rec.recommend(user, ALGO, 10)
    print(result)


if __name__ == "__main__":
    main()
