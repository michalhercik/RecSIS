import pandas as pd
import torch
from data import TrainData
from ranker.ranker import Ranker, cached
from torch_geometric.nn.models import LightGCN as TorchLightGCN
from torch_geometric.utils import negative_sampling
from user import User

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


class LightGCNRanker(Ranker):
    def fit(self):
        epochs = 10
        lr = 1e-2

        super().fit()

        def increment(df):
            df["course_id"] = df["course_id"] + self.train_data.user.shape[0]

        increment(self.train_data.povinn)
        increment(self.train_data.train)
        increment(self.train_data.val)

        # self.povinn["course_id"] = self.povinn["course_id"] + self.user.shape[0]
        # self.train["course_id"] = self.train["course_id"] + self.user.shape[0]
        # self.val["course_id"] = self.val["course_id"] + self.user.shape[0]

        self.id_to_povinn = dict(
            zip(self.train_data.povinn["course_id"], self.train_data.povinn["povinn"])
        )

        val_results = pd.merge(
            self.train_data.train.groupby("user_id")
            .agg({"course_id": set})
            .rename(columns={"course_id": "train_courses"}),
            self.train_data.val.groupby("user_id")
            .agg({"course_id": list})
            .rename(columns={"course_id": "val_courses"}),
            on="user_id",
        ).reset_index()

        num_nodes = self.train_data.user.shape[0] + self.train_data.povinn.shape[0]

        self.edge_index_homo = torch.stack(
            [
                torch.tensor(self.train_data.train["user_id"].values),
                torch.tensor(self.train_data.train["course_id"].values),
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

                pos_u, pos_i = self.edge_index_homo[:, : self.train_data.train.shape[0]]
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

    def rank(self, user: User) -> list[str]:
        uid = self.train_data.get_user(user.id)["user_id"].iloc[0]

        pred = self.model.recommend(
            edge_index=self.edge_index_homo,
            src_index=torch.tensor(uid),
            k=self.train_data.povinn.shape[0],
        )
        pred = pred.detach().cpu().tolist()
        pred = [self.id_to_povinn.get(cid) for cid in pred if cid in self.id_to_povinn]
        return pred
