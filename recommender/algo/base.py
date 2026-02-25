import pandas as pd

from data_repository import DataRepository
from user import User

class Student:

    def __init__(self, soident: str, degree_plan: str, year_of_study: int, type: str):
        self.soident = soident
        self.degree_plan = degree_plan
        self.year_of_study = year_of_study
        self.type = type

class RecommendedCourse:
    code: str
    true_positive: bool

    def __init__(self, code: str, true_positive: bool):
        self.code = code
        self.true_positive = true_positive

class ExpectedCourse:
    code: str
    false_negative: bool

    def __init__(self, code: str, false_negative: bool):
        self.code = code
        self.false_negative = false_negative

class Result:
    soident: str
    sobor: str
    degree_plan: str
    year_of_study: int
    type: str
    finished: list[str]
    finished_in_degree_plan: list[bool]
    recommended: list[str]
    recommended_true_positive: list[bool]
    recommended_in_degree_plan: list[bool]
    expected: list[str]
    expected_false_negative: list[bool]
    expected_in_degree_plan: list[bool]

    def __init__(self, soident: str, sobor: str, degree_plan: str, year_of_study: int, type: str, finished: list[str], finished_in_degree_plan: list[bool], recommended: list[RecommendedCourse] = None, recommended_in_degree_plan: list[bool] = None, expected: list[ExpectedCourse] = None, expected_in_degree_plan: list[bool] = None):
        self.soident = soident
        self.sobor = sobor
        self.degree_plan = degree_plan
        self.year_of_study = year_of_study
        self.type = type
        self.finished = finished
        self.finished_in_degree_plan = finished_in_degree_plan
        self.recommended = recommended
        self.recommended_in_degree_plan = recommended_in_degree_plan
        self.expected = expected
        self.expected_in_degree_plan = expected_in_degree_plan

    def recommended_and_expected(self, recommended: list[str], expected: list[str]):
        self.recommended = recommended
        self.recommended_true_positive = [code in expected for code in recommended]
        self.expected = expected
        self.expected_false_negative = [code not in recommended for code in expected]
        return self

class Recommendation:
    def __init__(self, rec: list[str], target: list[str], finished: list[str], expected: list[str]):
        self.rec = rec
        self.target = target
        self.finished = finished
        self.expected = expected

class Algorithm:
    def __init__(self, data: DataRepository):
        self.data = data

    def recommend(self, user: User, limit: int) -> Result:
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
