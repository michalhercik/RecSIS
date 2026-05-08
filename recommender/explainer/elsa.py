import pandas as pd
from data import TrainData
from explainer.explainer import Explainer
from ranker.elsa_ranker import Elsa
from user import User


class ElsaExplainer(Explainer):
    def __init__(self, model: Elsa, train_data: TrainData):
        super().__init__(train_data)
        self.model = model

    def fit(self) -> None:
        pass

    def explain(self, user: User, courses: list[str]) -> dict[str, str]:
        povinn = self.train_data.povinn
        order = {v: i for i, v in enumerate(courses)}
        courses = povinn[povinn["povinn"].isin(courses)].assign(
            _order=lambda df: df["povinn"].map(order)
        )
        finished = self.train_data.get_finished(user.id)
        finished = povinn[povinn["povinn"].isin(finished)]

        sim = self.model.similar_items(
            N=5,
            batch_size=3,
            sources=courses["course_id"].values,
            candidates=finished["course_id"].values,
        )[0].numpy()

        expl = (
            pd.DataFrame(sim, index=courses.index)
            .melt(ignore_index=False, value_name="sim", var_name="var_name")
            .reset_index()
            .merge(finished, left_on="sim", right_on="course_id")
            .merge(self.train_data.categories, on="povinn", how="left")
            .rename(columns={"povinn": "sim_povinn", "nazev": "sim_nazev"})
            .merge(courses, left_on="index", right_on="course_id")
            .merge(self.train_data.categories, on="povinn")
        )
        expl = expl[expl["nazev"] == expl["sim_nazev"]]
        expl = expl.groupby("index")["nazev"].agg(pd.Series.mode)

        expl = (
            courses[["_order", "povinn", "course_id", "pnazev", "pgarant"]]
            .merge(
                expl,
                left_index=True,
                right_index=True,
                how="left",
            )
            .sort_values("_order")
        )

        return dict(
            zip(
                expl["povinn"],
                expl["nazev"].astype(str),
            )
        )
