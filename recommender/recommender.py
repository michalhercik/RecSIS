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
from ranker.embedder import ContentKNN, UserKNN
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

    def explain(self, user: User, courses: list[str]) -> dict[str, str]:
        return self.explainer.explain(user, courses)


class EvalRecommender:
    def __init__(self):
        self.train_data = TrainData(rnd_state=RND_STATE)
        elsa = Elsa(self.train_data)
        self.model = {
            "Elsa": Model(elsa, ElsaExplainer(elsa, self.train_data)),
            "GCN": Model(GCNRanker(self.train_data), EmptyExplainer()),
            "LightGCN": Model(LightGCNRanker(self.train_data), EmptyExplainer()),
            "UserKNN": Model(UserKNN(self.train_data), EmptyExplainer()),
            "ContentKNN": Model(ContentKNN(self.train_data), EmptyExplainer()),
        }
        self.finished = FinishedFilter(self.train_data)
        self.grouper = SyntaxGrouper(self.train_data)
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

        expected = []
        finished = user.blueprint_to_df()["course"].to_list()
        if user.fetch:
            if user.id.lower() == "random":
                user.id = self.train_data.rand_soident_from_dev()

            expected = self.train_data.get_expected(user.id)
            finished = self.train_data.get_finished(user.id)

        self.degree_plan.fit(user)

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

        result = {
            "soident": str(user.id),
            "degree_plan": user.degree_plan,
            "finished": finished,
            "expected": expected,
            "recommended": recommended,
            "finished_in_degree_plan": self.degree_plan.mask(finished),
            "expected_in_degree_plan": self.degree_plan.mask(expected),
        }
        if user.fetch:
            result["type"] = self.train_data.get_type(user.id)
            result["sobor"] = self.train_data.get_sobor(user.id)
            result["degree_plan"] = self.train_data.get_degree_plan(user.id)
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
