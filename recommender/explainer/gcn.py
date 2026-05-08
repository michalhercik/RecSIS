import pandas as pd
import torch
from data import TrainData
from explainer.explainer import Explainer
from ranker.gcn import GCNRanker
from torch_geometric.explain import (
    CaptumExplainer,
    GNNExplainer,
    ModelConfig,
    PGExplainer,
)
from torch_geometric.explain import Explainer as TorchExplainer
from user import User


class ExplainerType:
    GNN = 1
    CAPTUM = 2


class Categories:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def ips(self) -> pd.DataFrame:
        # Compute inverse propensity scores (IPS) per category name
        total = len(self.train_data.categories)
        counts = self.train_data.categories["nazev"].value_counts()
        ips = total / counts

        # Apply log1p transformation (stabilize large values)
        vals = torch.log1p(torch.tensor(ips.values, dtype=torch.float32))

        # Min-max normalize to [0, 1]
        min_v = vals.min()
        max_v = vals.max()
        norm_vals = (vals - min_v) / (max_v - min_v)

        return pd.Series(norm_vals.numpy(), index=ips.index, name="ips")


class GCNExplainer(Explainer):
    def __init__(self, model: GCNRanker, train_data: TrainData):
        super().__init__(train_data)
        self.model = model
        self.categories = Categories(train_data)
        self.set_train_params()

    def fit(self) -> None:
        if self.type == ExplainerType.CAPTUM:
            self.explainer = self.Captum()
        elif self.type == ExplainerType.GNN:
            self.explainer = self.GNN()
        else:
            raise ValueError(f"Unknown explainer type: {self.type}")

    def GNN(self):
        return TorchExplainer(
            model=self.model.model,
            algorithm=GNNExplainer(epochs=self.epochs),
            explanation_type="model",
            model_config=ModelConfig(
                mode="binary_classification",
                task_level="edge",
                return_type="raw",
            ),
            node_mask_type="attributes",
            edge_mask_type="object",
        )

    def Captum(self):
        return TorchExplainer(
            model=self.model.model,
            algorithm=CaptumExplainer("IntegratedGradients"),
            explanation_type="model",
            model_config=ModelConfig(
                mode="binary_classification",
                task_level="edge",
                return_type="probs",
            ),
            node_mask_type="attributes",
            edge_mask_type="object",
            threshold_config=dict(
                threshold_type="topk",
                value=200,
            ),
        )

    def explain(self, user: User, courses: list[str]) -> list[str]:
        uid = self.train_data.user[self.train_data.user["soident"] == user.id]
        cids = self.train_data.povinn[self.train_data.povinn["povinn"].isin(courses)]
        uid = uid["user_id"].iloc[0]

        train_data = self.model.train_graph
        explanations = []
        for cid in cids["course_id"].tolist():
            c_tensor = torch.tensor([cid], dtype=torch.long)
            u_tensor = torch.tensor([uid], dtype=torch.long)
            edge_label_index = torch.stack([u_tensor, c_tensor], dim=0)
            exp = self.explainer(
                x=train_data.x_dict,
                edge_index=train_data.edge_index_dict,
                edge_label_index=edge_label_index,
            )
            explanations.append(exp)
        return self.node_level(user, explanations, cids)

    def node_level(self, user, explanations, cids):
        ips = self.categories.ips()
        label_expl = []
        for i, exp in enumerate(explanations):
            node_mask = exp["course"].node_mask
            mask_sum = node_mask.sum(dim=1)
            expl_course = pd.DataFrame(mask_sum, columns=["expl_score"])
            expl_course["expl"] = cids["course_id"].iloc[i]
            expl_course = expl_course.reset_index().rename(
                columns={"index": "course_id"}
            )
            expl_course = expl_course[
                (expl_course["expl_score"] > 0)
                & (expl_course["course_id"] != expl_course["expl"])
            ]
            expl_course["expl_score"] = (
                expl_course["expl_score"] / expl_course["expl_score"].sum()
            )
            expl_course = expl_course.merge(
                self.train_data.finished_df(user.id).drop_duplicates(), on="course_id"
            ).merge(self.train_data.categories, on="povinn", how="left")

            candidates = self.train_data.povinn[
                self.train_data.povinn["course_id"] == cids["course_id"].iloc[i]
            ].merge(self.train_data.categories, on="povinn")

            candidates = expl_course.merge(
                candidates, left_on="expl", right_on="course_id"
            )
            candidates = candidates[candidates["nazev_x"] == candidates["nazev_y"]]

            value_counts = candidates["nazev_x"].value_counts()
            expl_score = candidates.groupby("nazev_x").agg({"expl_score": "sum"})
            scores = (
                value_counts.to_frame()
                .join(expl_score, on="nazev_x")
                .join(ips, on="nazev_x")
            )
            scores["result"] = scores["ips"] * scores["expl_score"]
            if not scores.empty:
                expl = (scores["ips"] * scores["expl_score"]).idxmax()
            else:
                expl = pd.NA

            label_expl.append(expl)
        return label_expl

    def set_train_params(self, type=ExplainerType.GNN, epochs=20):
        self.epochs = epochs
        self.type = type
