import re

import numpy as np
import pandas as pd
from data import TrainData
from explainer.gcn import Categories

class Categorizer:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def categorize(self, courses: list[str]) -> tuple[list[str], list[list[str]]]:
        """
        Groups courses into categories.
        param:
            courses: list of course titles to group
            k: number of clusters to return (default: None, returns all clusters)
        return:
            list of category names and list of categories, where each category is a list of course titles,
            keeps order of first occurence of category in courses
        """
        raise NotImplementedError()

class DepartmentCategorizer(Categorizer):
    def categorize(self, courses: list[str]) -> tuple[list[str], list[list[str]]]:
        c = (
            pd.DataFrame({"povinn": courses, "rank": range(len(courses))})
            .merge(self.train_data.povinn[["povinn", "pgarant"]], on="povinn")
        )
        dep = (
            c[c["rank"] < 50]
            .groupby("pgarant", as_index=False)
            .agg({"rank": "mean"})
            .rename(columns={"rank": "dep_rank"})
        )
        c = (
            c.merge(dep, on="pgarant", how="left")
            .dropna(subset=["dep_rank"])
            .sort_values(["dep_rank", "rank"])
        )
        categories = (
            c.groupby("pgarant", as_index=False)
            .agg({"povinn": list})
        )
        return categories["pgarant"].to_list(), categories["povinn"].to_list()


class RankCategorizer(Categorizer):
    def categorize(self, courses: list[str]) -> tuple[list[str], list[list[str]]]:
        ips = Categories(self.train_data).ips()
        courses = courses[:50]
        pred = pd.DataFrame({"povinn": courses, "rank": range(len(courses))})
        pred = pred.merge(self.train_data.categories, on="povinn")
        pred["rank"] = 1 / np.log1p(pred["rank"] + 1)
        categories = pred.groupby("nazev", as_index=False).agg({"rank": "sum"})
        categories = categories.merge(ips, on="nazev")
        categories["ips"] = categories["ips"].pow(0.7)
        categories["result"] = categories["rank"] * categories["ips"]

        top_categories = (
            categories.sort_values(by="result", ascending=False)
            .head(5)["nazev"]
            .to_list()
        )
        cat_courses = [
            pred[pred["nazev"] == cat].sort_values(by="rank")["povinn"].to_list()
            for cat in top_categories
        ]

        return top_categories, cat_courses
