from data import TrainData
from explainer.elsa import ElsaExplainer
from explainer.explainer import EmptyExplainer, Explainer
from filterer import FinishedFilter
from grouper import IdentityGrouper, SyntaxGrouper
from categorizer import RankCategorizer, DepartmentCategorizer
from masker import (
    InMasker,
    NotInMasker,
)
from ranker.elsa_ranker import Elsa
from ranker.embedder import ContentKNN, UserKNN
from ranker.gcn import GCNRanker
from ranker.lightgcn import LightGCNRanker
from ranker.ranker import Ranker
from user import User

RND_STATE = 67564

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

class ProdRecommender:
    def __init__(self):
        self.train_data = TrainData(rnd_state=RND_STATE)
        elsa = Elsa(self.train_data)
        self.model = Model(elsa, ElsaExplainer(elsa, self.train_data))
        self.finished = FinishedFilter(self.train_data)
        self.grouper = SyntaxGrouper(self.train_data)
        self.categorizer = DepartmentCategorizer(self.train_data)


    def recommend(self, user: User, offset: int, limit: int, groups: bool, categories: bool):
        result = dict()
        # degree_plan = set(self.train_data.degree_plan_courses_by_code(user.degree_plan))
        finished = user.blueprint_to_df()["course"].to_list()
        ranking = self.model.rank(user)
        ranking = self.finished.filter(user, ranking)
        all = []

        if groups:
            result["groups"] = self.grouper.group(ranking[offset:], limit)
            limit = sum([len(group) for group in result["groups"]])
            result["pred"] = ranking[offset:offset+limit]
        else:
            result["pred"] = ranking[offset:offset+limit]

        all.extend(result["pred"])

        # explain = model.explain(user, pred)
        # result["explanations"] = explain

        if categories:
            cat_names, cat_values = self.categorizer.categorize(ranking)
            cat_values = [cat[:limit] for cat in cat_values]
            cat_groups = [self.grouper.group(c) for c in cat_values]
            result["categories"] = {
                "names": cat_names,
                "pred": cat_values,
                "groups": cat_groups,
            }

        for cat in result["categories"]["pred"]:
            all.extend(cat)
        result["all"] = [c for c in ranking if c in set(all)]

        return result

    def fit(self, cache=False):
        self.train_data.fit(cache=cache)
        self.grouper.fit()
        self.model.fit()


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
        self.notin_masker = NotInMasker()
        self.in_masker = InMasker()
        self.categorizer = RankCategorizer(self.train_data)

        self.train_data.fit(cache=True)
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
        degree_plan = set(self.train_data.degree_plan_courses_by_code(user.degree_plan))
        finished = user.blueprint_to_df()["course"].to_list()
        if user.fetch:
            if user.id.lower() == "random":
                user.id = self.train_data.rand_soident_from_dev()

            expected = self.train_data.get_expected(user.id)
            finished = self.train_data.get_finished(user.id)
            degree_plan = set(self.train_data.degree_plan_courses_by_soident(user.id))

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
                    "in_degree_plan": self.in_masker.mask(degree_plan, pred),
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
            "finished_in_degree_plan": self.in_masker.mask(degree_plan, finished),
            "expected_in_degree_plan": self.in_masker.mask(degree_plan, expected),
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
