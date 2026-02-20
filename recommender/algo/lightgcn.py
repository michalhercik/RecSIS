import numpy as np
import pandas as pd
import torch
from torch_geometric.nn.models import LightGCN
from torch_geometric.utils import negative_sampling

from algo.base import Algorithm
from algo.train import TrainData
from user import User

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# TODO: can recommned only for already seen data

# class LightGCN(Algorithm):
# def fit(self):
#     VAL_RATIO = 0.2
#     epochs = 1
#     lr = 1e-2

#     user, finished, povinn = self.dataset()
#     povinn["course_id"] = povinn["course_id"] + user.shape[0]
#     finished["course_id"] = finished["course_id"] + user.shape[0]
#     train, val, test = self.split(finished, VAL_RATIO, 2024)

#     val_results = pd.merge(
#         train.groupby("user_id")
#         .agg({"course_id": set})
#         .rename(columns={"course_id": "train_courses"}),
#         val.groupby("user_id")
#         .agg({"course_id": list})
#         .rename(columns={"course_id": "val_courses"}),
#         on="user_id",
#     ).reset_index()

#     num_nodes = user.shape[0] + povinn.shape[0]
#     self.model = LightGCN(num_nodes=num_nodes, embedding_dim=64, num_layers=3).to(device)

#     edge_index_homo = torch.stack(
#         [
#             torch.tensor(train["user_id"].values),
#             torch.tensor(train["course_id"].values),
#         ],
#         dim=0,
#     )
#     edge_index_homo = torch.cat([edge_index_homo, edge_index_homo.flip(0)], dim=1)

#     edge_index_homo = edge_index_homo
#     num_nodes = num_nodes
#     train = train

#     optimizer = torch.optim.AdamW(self.model.parameters(), lr=lr)

#     def train_step():
#         self.model.train()
#         optimizer.zero_grad()

#         neg_edge_index = negative_sampling(
#             edge_index=edge_index_homo,
#             num_nodes=num_nodes,
#             num_neg_samples=edge_index_homo.size(1) // 2,
#         )

#         pos_u, pos_i = edge_index_homo[:, : train.shape[0]]
#         _, neg_i = neg_edge_index

#         emb = self.model.get_embedding(edge_index_homo)

#         u_emb = emb[pos_u]
#         pos_emb = emb[pos_i]
#         neg_emb = emb[neg_i]

#         pos_scores = (u_emb * pos_emb).sum(dim=1)
#         neg_scores = (u_emb * neg_emb).sum(dim=1)

#         loss = self.model.recommendation_loss(
#             pos_scores,
#             neg_scores,
#             node_id=torch.cat([pos_u, pos_i, neg_i]),
#             lambda_reg=1e-4,
#         )
#         loss.backward()
#         optimizer.step()

#         return float(loss.detach())

#     for epoch in range(1, self.epochs + 1):
#         loss = train_step()
#         print(f"Epoch {epoch:03d} | Loss: {loss:.4f}")

#     self.model.eval()

# def recommend(self, user: User, limit: int) -> list[str]:
#     top_items = self.model.recommend(
#         edge_index=self.edge_index_homo,
#         src_index=torch.tensor(val_results["user_id"].values),
#         k=povinn.shape[0],
#     )
