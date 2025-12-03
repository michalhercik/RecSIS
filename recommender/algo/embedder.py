import os

import meilisearch
import numpy as np
import pandas as pd
from algo.base import Algorithm
from user import User


class Embedder(Algorithm):
    def fit(self):
        host = os.environ.get("MEILI_HOST", "http://meilisearch:7700")
        master_key = os.environ["MEILI_MASTER_KEY"]
        self.client = meilisearch.Client(host, master_key)
        self.course_index = self.client.index("courses")

    def recommend(self, user: User, limit: int) -> list[str]:
        bp = user.blueprint_to_df()
        query = self.build_query(bp)
        filter = self.build_filter(bp)
        similar = self.fetch_similar(query, filter, limit)
        return similar

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


class EmbedderAnnotation(Embedder):
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


class EmbedderSyllabus(Embedder):
    def build_query(
        self,
        blueprint_courses,
        prefix="Give me recommendations for similar courses like:",
    ):
        t = self.data.pamela
        annot = t[
            t["povinn"].isin(blueprint_courses)
            & (t["jazyk"] == "ENG")
            & (t["typ"] == "S")
        ]
        query = pd.merge(annot, self.data.povinn, how="left", on="povinn")
        query["panazev"] = query["panazev"].fillna("")
        query["memo"] = query["memo"].fillna("")
        query = query["panazev"] + ", " + query["memo"]
        query = "\n".join(query)
        return query
