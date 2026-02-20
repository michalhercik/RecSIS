import os

import numpy as np
import pandas as pd
import torch
from sentence_transformers import SentenceTransformer, models

from algo.train import TrainData
from user import User

RND_STATE = 42
VAL_RATIO = 0.2


# class UserKNN(TrainData):
#     def fit(self):
#         user, finished, povinn = self.dataset()
#         train, val, test = self.split(finished, VAL_RATIO, 2024)

#         pamela = self.data.pamela
#         pamela = pamela[
#             (pamela["jazyk"] == "ENG") & (pamela["typ"].isin(["A", "S"]))
#         ].pivot_table(index="povinn", columns="typ", values="memo", aggfunc="first")
#         povinn = povinn.merge(pamela, on="povinn")
#         embed_src = povinn.apply(
#             lambda x: f"{x['panazev']}: {x['A']}\n{x['S']}", axis=1
#         )
#         povinn["embed"] = list(sbert_embed(embed_src))

#         self.train_data = (
#             train.merge(povinn[["course_id", "embed"]], on="course_id")
#             .groupby("user_id")
#             .agg(
#                 {
#                     "user_id": "first",
#                     "course_id": set,
#                     "embed": lambda x: np.mean(x.values, axis=0),
#                 }
#             )
#             .rename(columns={"course_id": "train_courses"})
#             .reset_index(drop=True)
#         )

#         self.povinn = povinn

#     def recommend(self, user: User, limit: int) -> list[str]:
#         # TODO: transform user
#         dp = user.blueprint_to_df()
#         embed = dp.merge(self.povinn, left_on="course", right_on="povinn")
#         embed = embed["embed"].mean()
#         # pred = self.similar(self.povinn, embed)
#         pred = self.recommend_user(self.train_data.copy(), embed)
#         pred = self.filter_out_finished(user, pred)
#         pred = self.filter_out_dp(user, pred)
#         print(pred[:limit])
#         return pred[:limit]

#     def recommend_user(self, results, embed):
#         def cos_sim(x1, x2):
#             return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))

#         results["sim"] = results.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
#         results = results[["train_courses", "sim"]].explode("train_courses")
#         results = results.sort_values("sim", ascending=False)
#         results = results.drop_duplicates(subset="train_courses", keep="first")
#         results = results.merge(
#             self.povinn[["course_id", "povinn"]],
#             left_on="train_courses",
#             right_on="course_id",
#         )
#         return results["povinn"].to_list()

#     def similar(self, povinn, embed):
#         def cos_sim(x1, x2):
#             return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))

#         povinn["sim"] = povinn.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
#         povinn = povinn.sort_values("sim", ascending=False)
#         return povinn["povinn"].to_list()


# class ItemKNN(TrainData):
#     def fit(self):
#         user, finished, povinn = self.dataset()
#         train, val, test = self.split(finished, VAL_RATIO, 2024)

#         pamela = self.data.pamela
#         pamela = pamela[
#             (pamela["jazyk"] == "ENG") & (pamela["typ"].isin(["A", "S"]))
#         ].pivot_table(index="povinn", columns="typ", values="memo", aggfunc="first")
#         povinn = povinn.merge(pamela, on="povinn")
#         embed_src = povinn.apply(
#             lambda x: f"{x['panazev']}: {x['A']}\n{x['S']}", axis=1
#         )
#         povinn["embed"] = list(sbert_embed(embed_src))

#         self.train_data = (
#             train.merge(povinn[["course_id", "embed"]], on="course_id")
#             .groupby("user_id")
#             .agg(
#                 {
#                     "user_id": "first",
#                     "course_id": set,
#                     "embed": lambda x: np.mean(x.values, axis=0),
#                 }
#             )
#             .rename(columns={"course_id": "train_courses"})
#             .reset_index(drop=True)
#         )

#         self.povinn = povinn

#     def recommend(self, user: User, limit: int) -> list[str]:
#         # TODO: transform user
#         dp = user.blueprint_to_df()
#         embed = dp.merge(self.povinn, left_on="course", right_on="povinn")
#         embed = embed["embed"].mean()
#         # pred = self.similar(self.povinn, embed)
#         pred = self.recommend_user(self.train_data.copy(), embed)
#         pred = self.filter_out_finished(user, pred)
#         pred = self.filter_out_dp(user, pred)
#         print(pred[:limit])
#         return pred[:limit]

#     def recommend_user(self, results, embed):
#         def cos_sim(x1, x2):
#             return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))

#         results["sim"] = results.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
#         results = results[["train_courses", "sim"]].explode("train_courses")
#         results = results.sort_values("sim", ascending=False)
#         results = results.drop_duplicates(subset="train_courses", keep="first")
#         results = results.merge(
#             self.povinn[["course_id", "povinn"]],
#             left_on="train_courses",
#             right_on="course_id",
#         )
#         return results["povinn"].to_list()

#     def similar(self, povinn, embed):
#         def cos_sim(x1, x2):
#             return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))

#         povinn["sim"] = povinn.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
#         povinn = povinn.sort_values("sim", ascending=False)
#         return povinn["povinn"].to_list()


# consider Model2Vec https://github.com/MinishLab/model2vec?tab=readme-ov-file
word_embedding_model = models.Transformer("sentence-transformers/all-MiniLM-L12-v2")
pooling_model = models.Pooling(
    word_embedding_model.get_word_embedding_dimension(), pooling_mode_mean_tokens=True
)

sbert = SentenceTransformer(modules=[word_embedding_model, pooling_model])


def sbert_embed(texts):
    with torch.no_grad():
        embeddings = sbert.encode(texts, normalize_embeddings=True)
    return embeddings
