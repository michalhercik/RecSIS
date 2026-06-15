import re

import numpy as np
import pandas as pd
from data import TrainData
from explainer.gcn import Categories


class Grouper:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def fit(self) -> None:
        raise NotImplementedError()

    def group(self, courses: list[str], k=None) -> list[list[int]]:
        """
        Groups courses into clusters.
        param:
            courses: list of course titles to group
            k: number of clusters to return (default: None, returns all clusters)
        return:
            list of clusters, where each cluster is a list of course titles,
            keeps order of first occurence of cluster in courses
        """
        raise NotImplementedError()


class IdentityGrouper(Grouper):
    def __init__(self):
        pass

    def fit(self) -> None:
        pass

    def group(self, courses: list[str], k=None) -> list[list[int]]:
        if k is None:
            k = len(courses)
        return [[i] for i in range(k)]


class SyntaxGrouper(Grouper):
    roman_re: re.Pattern
    clusters: dict[str, int]
    cluster_names: dict[int, str]

    def fit(self) -> None:
        self.roman_re = re.compile(
            r"\b(?:M{0,4}(?:CM|CD|D?C{0,3})(?:XC|XL|L?X{0,3})(?:IX|IV|V?I{0,3}))\b",
            flags=re.IGNORECASE,
        )
        norm = [self.clean_title(c) for c in self.train_data.povinn["pnazev"]]
        vocab = {w: i for i, w in enumerate(set(norm))}
        clusters = [vocab[w] for w in norm]
        self.cluster_names = {cluster: name for name, cluster in vocab.items()}
        self.clusters = {
            course: cluster
            for course, cluster in zip(self.train_data.povinn["povinn"], clusters)
        }

    def group(self, courses: list[str], k=None) -> list[list[int]]:
        cluster_memory = {}
        result = []
        for i, course in enumerate(courses):
            if len(result) == k:
                break
            cluster = self.clusters[course]
            if cluster in cluster_memory:
                j = cluster_memory[cluster]
                result[j].append(i)
            else:
                result.append([i])
                cluster_memory[cluster] = len(result) - 1
        return result

    def clean_title(self, s):
        s = re.sub(r"\d+", "", s)  # remove Arabic digits
        s = self.roman_re.sub("", s)  # remove Roman numerals
        stop_words = set(
            [
                "pro začátečníky",
                "pro mírně pokročilé",
                "pro středně pokročilé",
                "pro pokročilé",
                "Pokročilé",
            ]
        )
        for w in stop_words:
            s = s.replace(w, "")
        s = re.sub(r"\s+", " ", s).strip()  # collapse whitespace
        return s


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
