import numpy as np
import pandas as pd
import torch
from algo.base import Result
from algo.embedder import UserKNN
from algo.train import cached
from torch_geometric.nn.models import LightGCN as TorchLightGCN
from torch_geometric.utils import negative_sampling
from user import User

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# TODO: can recommned only for already seen data


class LightGCN(UserKNN):
    def fit(self):
        VAL_RATIO = 0.2
        epochs = 10
        lr = 1e-2

        super().fit()

        self.povinn["course_id"] = self.povinn["course_id"] + self.user.shape[0]
        self.train["course_id"] = self.train["course_id"] + self.user.shape[0]
        self.val["course_id"] = self.val["course_id"] + self.user.shape[0]

        self.id_to_povinn = dict(zip(self.povinn["course_id"], self.povinn["povinn"]))

        val_results = pd.merge(
            self.train.groupby("user_id")
            .agg({"course_id": set})
            .rename(columns={"course_id": "train_courses"}),
            self.val.groupby("user_id")
            .agg({"course_id": list})
            .rename(columns={"course_id": "val_courses"}),
            on="user_id",
        ).reset_index()

        num_nodes = self.user.shape[0] + self.povinn.shape[0]

        self.edge_index_homo = torch.stack(
            [
                torch.tensor(self.train["user_id"].values),
                torch.tensor(self.train["course_id"].values),
            ],
            dim=0,
        )
        self.edge_index_homo = torch.cat(
            [self.edge_index_homo, self.edge_index_homo.flip(0)], dim=1
        )

        # num_nodes = num_nodes
        #

        def train_model():
            model = TorchLightGCN(
                num_nodes=num_nodes, embedding_dim=64, num_layers=3
            ).to(device)
            optimizer = torch.optim.AdamW(model.parameters(), lr=lr)

            def train_step():
                model.train()
                optimizer.zero_grad()

                neg_edge_index = negative_sampling(
                    edge_index=self.edge_index_homo,
                    num_nodes=num_nodes,
                    num_neg_samples=self.edge_index_homo.size(1) // 2,
                )

                pos_u, pos_i = self.edge_index_homo[:, : self.train.shape[0]]
                _, neg_i = neg_edge_index

                emb = model.get_embedding(self.edge_index_homo)

                u_emb = emb[pos_u]
                pos_emb = emb[pos_i]
                neg_emb = emb[neg_i]

                pos_scores = (u_emb * pos_emb).sum(dim=1)
                neg_scores = (u_emb * neg_emb).sum(dim=1)

                loss = model.recommendation_loss(
                    pos_scores,
                    neg_scores,
                    node_id=torch.cat([pos_u, pos_i, neg_i]),
                    lambda_reg=1e-4,
                )
                loss.backward()
                optimizer.step()

                return float(loss.detach())

            for epoch in range(1, epochs + 1):
                loss = train_step()
                print(f"Epoch {epoch:03d} | Loss: {loss:.4f}")

            model.eval()
            return model

        self.model = cached(train_model, "lightgcn.pickle")

    def recommend(self, user: User, limit: int) -> list[str]:
        result = Result()
        result.soident = self.get_user_soident(user)
        result.type, result.sobor, result.degree_plan = self.get_user_info(
            user, result.soident
        )
        result.year_of_study, result.finished = self.get_year_finished(
            user, result.soident
        )
        result.expected = self.get_expected(user, result.soident)

        uid = self.get_user_id(user, result.soident)

        pred = self.model.recommend(
            edge_index=self.edge_index_homo,
            src_index=torch.tensor(uid),
            k=self.povinn.shape[0],
        )
        pred = pred.detach().cpu().tolist()
        pred = [self.id_to_povinn.get(cid) for cid in pred if cid in self.id_to_povinn]
        pred = self.filter_out_finished(pred, result.finished)

        result.recommended = pred[:limit]
        dp_courses = self.get_degree_plan_courses(result.degree_plan)
        result.generate_masks(dp_courses)
        return result
