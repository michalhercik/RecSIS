import pandas as pd
import torch
from data import TrainData
from ranker.ranker import Ranker, cached
from torch_geometric.nn.models import LightGCN as TorchLightGCN
from torch_geometric.utils import negative_sampling
from user import User

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


class LightGCNRanker(Ranker):
    def __init__(self, train_data: TrainData):
        super().__init__(train_data)
        self.set_train_params(
            embedding_dim=64, num_layers=3, num_epochs=3, learning_rate=1e-2
        )

    def fit(self):
        # epochs = 3
        # lr = 1e-2

        self.edge_index_homo, self.id_to_povinn = self.edge_index(
            self.train_data.user.copy(),
            self.train_data.train.copy(),
            self.train_data.povinn.copy(),
        )
        num_nodes = self.train_data.user.shape[0] + self.train_data.povinn.shape[0]
        self.model = cached(
            lambda: self.fit_with(
                num_nodes, self.edge_index_homo, self.train_data.train.shape[0]
            ),
            "lightgcn.pickle",
        )

        # def increment(df):
        #     df["course_id"] = df["course_id"] + self.train_data.user.shape[0]

        # increment(self.train_data.povinn)
        # increment(self.train_data.train)
        # increment(self.train_data.val)

        # self.id_to_povinn = dict(
        #     zip(self.train_data.povinn["course_id"], self.train_data.povinn["povinn"])
        # )

        # num_nodes = self.train_data.user.shape[0] + self.train_data.povinn.shape[0]

        # self.edge_index_homo = torch.stack(
        #     [
        #         torch.tensor(self.train_data.train["user_id"].values),
        #         torch.tensor(self.train_data.train["course_id"].values),
        #     ],
        #     dim=0,
        # )
        # self.edge_index_homo = torch.cat(
        #     [self.edge_index_homo, self.edge_index_homo.flip(0)], dim=1
        # )

        # def train_model():
        #     model = TorchLightGCN(
        #         num_nodes=num_nodes, embedding_dim=64, num_layers=3
        #     ).to(device)
        #     optimizer = torch.optim.AdamW(model.parameters(), lr=lr)

        #     def train_step():
        #         model.train()
        #         optimizer.zero_grad()

        #         neg_edge_index = negative_sampling(
        #             edge_index=self.edge_index_homo,
        #             num_nodes=num_nodes,
        #             num_neg_samples=self.edge_index_homo.size(1) // 2,
        #         )

        #         pos_u, pos_i = self.edge_index_homo[:, : self.train_data.train.shape[0]]
        #         _, neg_i = neg_edge_index

        #         emb = model.get_embedding(self.edge_index_homo)

        #         u_emb = emb[pos_u]
        #         pos_emb = emb[pos_i]
        #         neg_emb = emb[neg_i]

        #         pos_scores = (u_emb * pos_emb).sum(dim=1)
        #         neg_scores = (u_emb * neg_emb).sum(dim=1)

        #         loss = model.recommendation_loss(
        #             pos_scores,
        #             neg_scores,
        #             node_id=torch.cat([pos_u, pos_i, neg_i]),
        #             lambda_reg=1e-4,
        #         )
        #         loss.backward()
        #         optimizer.step()

        #         return float(loss.detach())

        #     for epoch in range(1, epochs + 1):
        #         loss = train_step()
        #         print(f"Epoch {epoch:03d} | Loss: {loss:.4f}")

        #     model.eval()
        #     return model

        # self.model = cached(train_model, "lightgcn.pickle")

    def rank(self, user: User) -> list[str]:
        u = self.train_data.get_user(user.id)
        if not u.empty:
            uid = u["user_id"].iloc[0]
            model = self.model
            id_to_povinn = self.id_to_povinn
            edge_index = self.edge_index_homo
        else:
            blueprint_df = user.blueprint_to_df()
            if blueprint_df.empty:
                return []

            edge_index_homo = self.edge_index_homo
            emb = self.model.get_embedding(edge_index_homo).detach()
            course_id = (
                self.train_data.povinn["course_id"] + self.train_data.user.shape[0]
            )
            course_emb = emb[course_id]
            blueprint = blueprint_df.merge(
                self.train_data.povinn, left_on="course", right_on="povinn"
            )["course_id"]
            u_emb = emb[blueprint].mean(dim=0)
            results = self.train_data.povinn.copy()

            results["score"] = torch.matmul(
                torch.tensor(course_emb, dtype=torch.float32),
                torch.tensor(u_emb, dtype=torch.float32),
            )
            results = results.sort_values("score", ascending=False)
            pred = results["povinn"].tolist()
            return pred

        # else:
        #     uid, model, edge_index, id_to_povinn = self.train_with_user(user)

        pred = model.recommend(
            edge_index=edge_index,
            src_index=torch.tensor(uid),
            k=self.train_data.povinn.shape[0],
        )
        pred = pred.detach().cpu().tolist()
        pred = [id_to_povinn.get(cid) for cid in pred if cid in self.id_to_povinn]
        return pred

    def set_train_params(self, embedding_dim, num_layers, num_epochs, learning_rate):
        self.embedding_dim = embedding_dim
        self.num_layers = num_layers
        self.num_epochs = num_epochs
        self.learning_rate = learning_rate

    def train_with_user(
        self, user: User
    ) -> tuple[int, torch.nn.Module, torch.Tensor, dict[int, int]]:
        uid = self.train_data.user.shape[0]
        user_df = pd.concat(
            [
                self.train_data.user[["user_id", "soident"]],
                pd.DataFrame(
                    [
                        [
                            uid,
                            user.id,
                        ]
                    ],
                    index=[uid],
                    columns=["user_id", "soident"],
                ),
            ]
        )
        finished = user.blueprint_to_df()
        finished["user_id"] = uid
        finished = finished.merge(
            self.train_data.povinn[["course_id", "povinn"]],
            left_on="course",
            right_on="povinn",
            how="inner",
        )
        interactions = pd.concat(
            [self.train_data.train, finished[["user_id", "course_id"]]]
        )

        edge_index, id_to_povinn = self.edge_index(
            user_df, interactions, self.train_data.povinn
        )
        model = self.fit_with(
            user_df.shape[0] + self.train_data.povinn.shape[0],
            edge_index,
            interactions.shape[0],
        )
        return uid, model, edge_index, id_to_povinn

    def edge_index(
        self, user: pd.DataFrame, train: pd.DataFrame, povinn: pd.DataFrame
    ) -> tuple[torch.nn.Module, dict[int, int]]:
        def increment(df):
            df["course_id"] = df["course_id"] + user.shape[0]

        increment(povinn)
        increment(train)
        id_to_povinn = dict(zip(povinn["course_id"], povinn["povinn"]))

        edge_index_homo = torch.stack(
            [
                torch.tensor(train["user_id"].values),
                torch.tensor(train["course_id"].values),
            ],
            dim=0,
        )
        edge_index_homo = torch.cat([edge_index_homo, edge_index_homo.flip(0)], dim=1)
        return edge_index_homo, id_to_povinn

    def fit_with(
        self, num_nodes: int, edge_index_homo: torch.Tensor, num_edges: int
    ) -> torch.nn.Module:
        def train_model():
            model = TorchLightGCN(
                num_nodes=num_nodes,
                embedding_dim=self.embedding_dim,
                num_layers=self.num_layers,
            ).to(device)
            optimizer = torch.optim.AdamW(model.parameters(), lr=self.learning_rate)

            def train_step():
                model.train()
                optimizer.zero_grad()

                neg_edge_index = negative_sampling(
                    edge_index=edge_index_homo,
                    num_nodes=num_nodes,
                    num_neg_samples=edge_index_homo.size(1) // 2,
                )

                pos_u, pos_i = edge_index_homo[:, :num_edges]
                _, neg_i = neg_edge_index

                emb = model.get_embedding(edge_index_homo)

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

            for epoch in range(1, self.num_epochs + 1):
                loss = train_step()
                print(f"Epoch {epoch:03d} | Loss: {loss:.4f}")

            model.eval()
            return model

        return train_model()
