import os
import pickle

import meilisearch
import numpy as np
import pandas as pd
import torch
from algo.base import Algorithm, Result
from algo.train import TrainData, cached
from sentence_transformers import SentenceTransformer, models
from transformers import (
    BertModel,
    BertTokenizer,
)
from user import User

RND_STATE = 42
VAL_RATIO = 0.2


def cos_sim(x1, x2):
    return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))


class UserKNN(TrainData):
    def fit(self):
        super().fit()

        povinn_file = "embedder-povinn.pickle"
        train_data_file = "embedder-train-data.pickle"

        def embed_povinn():
            pamela = self.data.pamela
            pamela = pamela[
                (pamela["jazyk"] == "ENG") & (pamela["typ"].isin(["A", "S"]))
            ].pivot_table(index="povinn", columns="typ", values="memo", aggfunc="first")
            embed_povinn = self.povinn.merge(pamela, on="povinn")
            embed_src = embed_povinn.apply(
                lambda x: f"{x['panazev']}: {x['A']}\n{x['S']}", axis=1
            )
            embed_povinn["embed"] = list(sbert_embed(embed_src))
            return embed_povinn

        def train_data():
            result = (
                self.train.merge(
                    self.embed_povinn[["course_id", "embed"]], on="course_id"
                )
                .groupby("user_id")
                .agg(
                    {
                        "user_id": "first",
                        "course_id": set,
                        "embed": lambda x: np.mean(x.values, axis=0),
                    }
                )
                .rename(columns={"course_id": "train_courses"})
                .reset_index(drop=True)
            )
            return result

        self.embed_povinn = cached(
            embed_povinn,
            povinn_file,
            lambda: os.path.exists(povinn_file) and os.path.exists(train_data_file),
        )

        self.train_data = cached(train_data, train_data_file)

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

        user_embedding = self.embed(user, result.soident)
        pred = self.similar_users_finished_courses(user_embedding)
        pred = self.filter_out_finished(pred, result.finished)
        # pred = self.filter_out_dp(user, pred)
        pred = pred[:limit]
        result.recommended = pred

        dp_courses = self.get_degree_plan_courses(result.degree_plan)
        result.generate_masks(dp_courses)
        return result

    def similar_users_finished_courses(self, embed):
        results = self.train_data.copy()
        results["sim"] = results.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
        results = results[["train_courses", "sim"]].explode("train_courses")
        results = results.sort_values("sim", ascending=False)
        results = results.drop_duplicates(subset="train_courses", keep="first")
        results = results.merge(
            self.povinn[["course_id", "povinn"]],
            left_on="train_courses",
            right_on="course_id",
        )
        results = results["povinn"].to_list()
        return results

    def embed(self, user: User, soident: str):
        embed = None
        if user.fetch:
            user_row = self.user[self.user["soident"] == soident]
            finished = user_row.merge(self.val, on="user_id")
            if finished.empty:
                all = self.user.merge(self.train_data, on="user_id")
                embed = all.sample(1, random_state=RND_STATE)["embed"].iloc[0]
                return embed
            embed = finished.merge(self.embed_povinn, on="course_id")
        else:
            finished = user.blueprint_to_df()
            embed = finished.merge(
                self.embed_povinn, left_on="course", right_on="povinn"
            )

        embed = embed["embed"].mean()
        return embed


class ItemKNN(UserKNN):
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

        user_embedding = self.embed(user, result.soident)
        pred = self.similar(user_embedding)
        pred = self.filter_out_finished(pred, result.finished)
        # pred = self.filter_out_dp(user, pred)
        pred = pred[:limit]
        result.recommended = pred

        dp_courses = self.get_degree_plan_courses(result.degree_plan)
        result.generate_masks(dp_courses)
        return result

    def similar(self, embed):
        povinn = self.embed_povinn[["povinn", "embed"]].copy()
        povinn["sim"] = povinn.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
        povinn = povinn.sort_values("sim", ascending=False)
        return povinn["povinn"].to_list()


class MeiliSearchKNN(TrainData):
    def fit(self):
        super().fit()
        host = os.environ.get("MEILI_HOST", "http://localhost:7700")
        master_key = os.environ["MEILI_MASTER_KEY"]
        self.client = meilisearch.Client(host, master_key)
        self.course_index = self.client.index("courses")

    def recommend(self, user: User, limit: int) -> Result:
        raise NotImplementedError()

        bp = user.blueprint_to_df()
        query = self.build_query(bp)
        filter = self.build_filter(bp)
        similar = self.fetch_similar(query, filter, limit)
        return Result()

    def build_query(
        self,
        blueprint_courses,
        prefix="Give me recommendations for similar courses like:",
    ):
        povinn = self.data.povinn.loc[
            self.data.povinn["povinn"].isin(blueprint_courses)
        ]
        povinn_str = ",".join(povinn["panazev"])
        query = prefix + povinn_str
        return query

    def build_filter(self, course_codes):
        filter = "code NOT IN ['" + "','".join(course_codes) + "']"
        filter += " AND section=NI"
        return filter

    def fetch_similar(self, query, filter, limit):
        result = self.course_index.search(
            query,
            {
                "hybrid": {"semanticRatio": 1, "embedder": "bert"},
                "filter": filter,
                "attributesToRetrieve": ["code"],
                "limit": limit,
            },
        )
        codes = map(lambda x: x["code"], result["hits"])
        return list(codes)


class MeiliSearchKNNWithAnnotation(MeiliSearchKNN):
    def build_query(
        self,
        blueprint_courses,
        prefix="Give me recommendations for similar courses like:",
    ):
        t = self.data.pamela
        annot = t[
            t["povinn"].isin(blueprint_courses)
            & (t["jazyk"] == "ENG")
            & (t["typ"] == "A")
        ]
        query = pd.merge(annot, self.data.povinn, how="left", on="povinn")
        query["panazev"] = query["panazev"].fillna("")
        query["memo"] = query["memo"].fillna("")
        query = query["panazev"] + ", " + query["memo"]
        query = "\n".join(query)
        return query


# class EmbedderSyllabus(Embedder):
#     def build_query(
#         self,
#         blueprint_courses,
#         prefix="Give me recommendations for similar courses like:",
#     ):
#         t = self.data.pamela
#         annot = t[
#             t["povinn"].isin(blueprint_courses)
#             & (t["jazyk"] == "ENG")
#             & (t["typ"] == "S")
#         ]
#         query = pd.merge(annot, self.data.povinn, how="left", on="povinn")
#         query["panazev"] = query["panazev"].fillna("")
#         query["memo"] = query["memo"].fillna("")
#         query = query["panazev"] + ", " + query["memo"]
#         query = "\n".join(query)
#         return query


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
