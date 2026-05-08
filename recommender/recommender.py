from data import TrainData
from explainer.elsa import ElsaExplainer
from explainer.explainer import EmptyExplainer, Explainer
from filterer import FinishedFilter
from grouper import IdentityGrouper, RankCategorizer, SyntaxGrouper
from masker import (
    DegreePlanMasker,
    InMasker,
    NotInMasker,
    TruePositivesMasker,
)
from ranker.elsa_ranker import Elsa
from ranker.embedder import MeiliSearch
from ranker.gcn import GCNRanker
from ranker.lightgcn import LightGCNRanker
from ranker.ranker import Ranker
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


class Model(Ranker, Explainer):
    def __init__(self, ranker: Ranker, explainer: Explainer):
        self.ranker = ranker
        self.explainer = explainer

    def fit(self) -> None:
        self.ranker.fit()
        self.explainer.fit()

    def rank(self, user: User) -> list[str]:
        return self.ranker.rank(user)

    def explain(self, user: User, courses: list[str]) -> list[str]:
        return self.explainer.explain(user, courses)


class ModelFactory:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def elsa(self):
        ranker = Elsa(self.train_data)
        explainer = ElsaExplainer(ranker, self.train_data)
        return Model(ranker, explainer)

    def gcn(self):
        ranker = GCNRanker(self.train_data)
        explainer = EmptyExplainer()
        return Model(ranker, explainer)

    def light_gcn(self):
        ranker = LightGCNRanker(self.train_data)
        explainer = EmptyExplainer()
        return Model(ranker, explainer)


class EvalRecommender:
    def __init__(self):
        self.train_data = TrainData(rnd_state=RND_STATE)
        modelFactory = ModelFactory(self.train_data)
        self.model = {
            "Elsa": modelFactory.elsa(),
            "GCN": modelFactory.gcn(),
            "LightGCN": modelFactory.light_gcn(),
        }
        self.finished = FinishedFilter(self.train_data)
        self.grouper = SyntaxGrouper(self.train_data)
        self.true_pos = TruePositivesMasker(self.train_data)
        self.degree_plan = DegreePlanMasker(self.train_data)
        self.notin_masker = NotInMasker()
        self.in_masker = InMasker()
        self.categorizer = RankCategorizer(self.train_data)

        self.train_data.fit()
        self.grouper.fit()

        self.ranker_fitted = set()

    def algorithms(self):
        algos = list(self.model.keys())
        fitted = [True if algo in self.ranker_fitted else False for algo in algos]
        return {"algorithms": algos, "fit": fitted}

    def recommend(self, user: User, algos: list[str], limit: int):
        if not self.__is_known_algo(algos):
            return

        if user.id.lower() == "random":
            user.id = self.train_data.rand_soident_from_dev()

        self.true_pos.fit(user)
        self.degree_plan.fit(user)

        expected = self.train_data.get_expected(user.id)

        recommended = []
        for algo in algos:
            model = self.model[algo]
            ranking = model.rank(user)
            ranking = self.finished.filter(user, ranking)
            groups = self.grouper.group(ranking, limit)
            limit = sum([len(group) for group in groups])
            pred = ranking[:limit]
            explain = model.explain(user, pred)
            cat_names, cat_values = self.categorizer.categorize(ranking)
            cat_groups = [self.grouper.group(c) for c in cat_values]
            recommended.append(
                {
                    "pred": pred,
                    "groups": groups,
                    "true_positive": self.in_masker.mask(expected, pred),
                    "in_degree_plan": self.degree_plan.mask(pred),
                    "false_negative": self.notin_masker.mask(pred, expected),
                    "categories": {
                        "names": cat_names,
                        "pred": cat_values,
                        "groups": cat_groups,
                    },
                    "explanations": explain,
                }
            )

        finished = self.train_data.get_finished(user.id)
        result = {
            "soident": str(user.id),
            "type": self.train_data.get_type(user.id),
            "sobor": self.train_data.get_sobor(user.id),
            "degree_plan": self.train_data.get_degree_plan(user.id),
            "finished": finished,
            "expected": expected,
            "recommended": recommended,
            "finished_in_degree_plan": self.degree_plan.mask(finished),
            "expected_in_degree_plan": self.degree_plan.mask(expected),
        }
        # import json

        # print(json.dumps(result, indent=4), flush=True)
        return result

    def fit(self, algos: list[str]):
        if not self.__is_known_algo(algos):
            return
        for algo in algos:
            self.model[algo].fit()
            self.ranker_fitted.add(algo)

    def __is_known_algo(self, algos: list[str]) -> bool:
        return len(self.__unknown_algos(algos)) == 0

    def __unknown_algos(self, algos: list[str]) -> list[str]:
        return [algo for algo in algos if algo not in self.model]


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
