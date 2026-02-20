import pandas as pd

from data_repository import DataRepository
from user import User

class Recommendation:
    def __init__(self, rec: list[str], target: list[str], finished: list[str], expected: list[str]):
        self.rec = rec
        self.target = target
        self.finished = finished
        self.expected = expected

class Algorithm:
    def __init__(self, data: DataRepository):
        self.data = data

    def recommend(self, user: User, limit: int) -> Recommendation:
        raise NotImplementedError("Please Implement this method")

    def fit(self):
        raise NotImplementedError("Algorithm.fit() is not implemented")

    def filter_out_finished(self, user: User, predictions: list[str]) -> list[str]:
        finished = set(user.blueprint_to_df()["course"].to_list())
        predictions = [i for i in predictions if i not in finished]
        return predictions

    def filter_out_dp(self, user: User, predictions: list[str]) -> list[str]:
        dp = self.data.stud_plan
        dp = dp[dp["plan_code"] == user.degree_plan]
        dp = set(dp["code"].unique().tolist())
        predictions = [i for i in predictions if i not in dp]
        return predictions
