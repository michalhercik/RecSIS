from data import TrainData
from filterer import FinishedFilter
from grouper import IdentityGrouper, SyntaxGrouper
from masker import (
    DegreePlanMasker,
    FalseNegativeMasker,
    TruePositivesMasker,
)
from ranker.elsa import Elsa
from ranker.embedder import MeiliSearch
from user import User

RND_STATE = 67564


class Recommender:
    def __init__(self):
        self.train_data = TrainData(rnd_state=RND_STATE)
        self.ranker = Elsa(self.train_data)
        self.finished = FinishedFilter(self.train_data)
        self.grouper = SyntaxGrouper(self.train_data)

        self.train_data.fit()
        self.grouper.fit()

    def recommend(self, user: User, limit: int):
        ranking = self.ranker.rank(user)
        ranking = self.finished.filter(user, ranking)
        groups = self.grouper.group(ranking, limit)
        limit = sum([len(group) for group in groups])

        result = {
            "pred": ranking[:limit],
            "groups": groups,
        }
        return result

    def fit(self):
        self.ranker.fit()


class EvalRecommender:
    def __init__(self):
        self.train_data = TrainData(rnd_state=RND_STATE)
        self.ranker = {
            "Elsa": Elsa(self.train_data),
            # "MeiliSearch": MeiliSearch(self.train_data),
        }
        self.finished = FinishedFilter(self.train_data)
        self.grouper = SyntaxGrouper(self.train_data)
        self.true_pos = TruePositivesMasker(self.train_data)
        self.degree_plan = DegreePlanMasker(self.train_data)
        self.false_neg = FalseNegativeMasker(self.train_data)

        self.train_data.fit()
        self.grouper.fit()

        self.ranker_fitted = set()

    def algorithms(self):
        algos = list(self.ranker.keys())
        fitted = [True if algo in self.ranker_fitted else False for algo in algos]
        return {"algorithms": algos, "fit": fitted}

    def recommend(self, user: User, algos: list[str], limit: int):
        if not self.__is_known_algo(algos):
            return

        if user.id.lower() == "random":
            user.id = self.train_data.rand_soident_from_dev()

        self.true_pos.fit(user)
        self.degree_plan.fit(user)
        self.false_neg.fit(user)

        recommended = []
        for algo in algos:
            ranking = self.ranker[algo].rank(user)
            ranking = self.finished.filter(user, ranking)
            groups = self.grouper.group(ranking, limit)
            limit = sum([len(group) for group in groups])
            pred = ranking[:limit]
            recommended.append(
                {
                    "pred": pred,
                    "groups": groups,
                    "true_positive": self.true_pos.mask(pred),
                    "in_degree_plan": self.degree_plan.mask(pred),
                    "false_negative": self.false_neg.mask(pred),
                }
            )

        finished = self.train_data.get_finished(user.id)
        expected = self.train_data.get_expected(user.id)
        result = {
            "soident": user.id,
            "type": self.train_data.get_type(user.id),
            "sobor": self.train_data.get_sobor(user.id),
            "degree_plan": self.train_data.get_degree_plan(user.id),
            "finished": finished,
            "expected": expected,
            "recommended": recommended,
            "finished_in_degree_plan": self.degree_plan.mask(finished),
            "expected_in_degree_plan": self.degree_plan.mask(expected),
        }
        import json

        print(json.dumps(result, indent=4), flush=True)
        return result

    def fit(self, algos: list[str]):
        if not self.__is_known_algo(algos):
            return
        for algo in algos:
            self.ranker[algo].fit()
            self.ranker_fitted.add(algo)

    def __is_known_algo(self, algos: list[str]) -> bool:
        return len(self.__unknown_algos(algos)) == 0

    def __unknown_algos(self, algos: list[str]) -> list[str]:
        return [algo for algo in algos if algo not in self.ranker]


# class RecStudentInfo:
#     soident: str
#     type: str
#     sobor: str
#     degree_plan: str
#     year_of_study: int
#     finished: list[str]
#     expected: list[str]


# class GeneralInfo:
#     counter: int = 0

#     def __init__(self, train_data: TrainData):
#         self.train_data = train_data

#     def get(self, user: User) -> dict:
#         result = {}
#         result["soident"] = user.id
#         result["type"], result["sobor"], result["degree_plan"] = self.get_user_info(
#             user, user.id
#         )
#         result["year_of_study"] = self.get_year_finished(user)
#         return result

#     # def get_user_soident(self, user: User):
#     #     soident = user.id
#     #     if user.fetch:
#     #         if user.id.lower() == "random":
#     #             uid = (
#     #                 self.train_data.val["user_id"]
#     #                 .drop_duplicates()
#     #                 .sample(1, random_state=RND_STATE + self.counter)
#     #                 .iloc[0]
#     #             )
#     #             soident = self.train_data.user[self.train_data.user["user_id"] == uid][
#     #                 "soident"
#     #             ].iloc[0]
#     #             self.counter += 1
#     #     else:
#     #         soident = ""
#     #     return soident

#     def get_user_info(self, user: User, soident: str):
#         type = ""
#         sobor = ""
#         degree_plan = user.degree_plan
#         if user.fetch:
#             df = self.train_data.user[self.train_data.user["soident"] == int(soident)]
#             type = df["sdruh"].iloc[0]
#             sobor = df["sobor_nazev"].iloc[0]
#             degree_plan = df["splan"].iloc[0]
#         return type, sobor, degree_plan

#     def get_year_finished(self, user: User):
#         year = -1
#         if user.fetch:
#             user_id = self.train_data.user[
#                 self.train_data.user["soident"] == int(user.id)
#             ]["user_id"].iloc[0]
#             finished_df = self.train_data.val[self.train_data.val["user_id"] == user_id]
#             finished_df = finished_df.merge(self.train_data.povinn, on="course_id")
#             year = finished_df["zroc"].max()

#         return year

#     def get_expected(self, user: User, soident: str):
#         expected = []
#         if user.fetch:
#             user_id = self.train_data.user[
#                 self.train_data.user["soident"] == int(soident)
#             ]["user_id"].iloc[0]
#             expected_df = self.train_data.val[self.train_data.val["user_id"] == user_id]
#             expected_df = expected_df.merge(self.train_data.povinn, on="course_id")
#             expected = expected_df["povinn"].to_list()

#         return expected
