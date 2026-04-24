import os
import pickle

# import re
import sys

sys.path.append("..")
import pandas as pd
import psycopg2
from algo.base import Algorithm
from evaluation.tokenizer import equality_labels, group_by_cluster

# from scipy.cluster.hierarchy import fcluster, linkage
# from scipy.spatial.distance import squareform
# from sklearn.feature_extraction.text import CountVectorizer
# from sklearn.metrics import pairwise_distances
from user import User

RND_STATE = 42
VAL_RATIO = 0.2


def main():
    with open("povinn.pickle", "rb") as f:
        povinn = pickle.load(f)
    povinn["cluster"], cluster_names = equality_labels(povinn["pnazev"])
    # print(povinn)
    # povinn["cluster"], _, _, _ = assign_cluster_labels(
    #     povinn["pnazev"], distance_threshold=0.4
    # )
    # print(sorted(cluster_names.values()))
    # print(
    #     povinn.merge(povinn["cluster"].value_counts(), on="cluster")
    #     .sort_values(by=["count", "cluster"], ascending=False)[
    #         ["povinn", "pnazev", "cluster", "count"]
    #     ]
    #     .head(50)
    # )
    # print(povinn[povinn["cluster"] == 206]["pnazev"].tolist())
    # print(povinn[povinn["pnazev"].str.contains("výchova")])
    #
    print(povinn)
    pred = povinn[povinn["povinn"].isin(["NAIL069", "NAIL070"])]["course_id"].sample(
        frac=1
    )
    bla = povinn.iloc[pred]
    print(bla)
    cluster_memory = {}
    result = []
    k = 100
    for _, row in bla.iterrows():
        if len(result) == k:
            break
        cluster = row["cluster"]
        if cluster in cluster_memory:
            j = cluster_memory[cluster]
            result[j].append(row["povinn"])
        else:
            result.append([row["povinn"]])
            cluster_memory[cluster] = len(result) - 1
    print(result[:10])
    print([r for r in result if len(r) > 1])


def cached(retrieve, file, condition=None):
    if condition is None or condition():
        if os.path.exists(file):
            with open(file, "rb") as f:
                return pickle.load(f)
    data = retrieve()
    with open(file, "wb") as f:
        pickle.dump(data, f, protocol=pickle.HIGHEST_PROTOCOL)
    return data


class TrainData(Algorithm):
    def fit(self):
        self.counter = 0
        if (
            os.path.exists("povinn.pickle")
            and os.path.exists("val.pickle")
            and os.path.exists("user.pickle")
            and os.path.exists("train.pickle")
            and os.path.exists("finished.pickle")
            and os.path.exists("cluster_vocab.pickle")
        ):
            with open("povinn.pickle", "rb") as f:
                self.povinn = pickle.load(f)
            with open("val.pickle", "rb") as f:
                self.val = pickle.load(f)
            with open("user.pickle", "rb") as f:
                self.user = pickle.load(f)
            with open("train.pickle", "rb") as f:
                self.train = pickle.load(f)
            with open("finished.pickle", "rb") as f:
                self.finished = pickle.load(f)
            with open("cluster_vocab.pickle", "rb") as f:
                self.cluster_vocab = pickle.load(f)
            return

        user, finished, povinn = self.dataset()
        train, val, test = self.split(finished, VAL_RATIO, 2024)

        povinn["cluster"], vocab = equality_labels(povinn["pnazev"])

        self.user = user
        self.povinn = povinn
        self.val = val
        self.train = train
        self.finished = finished
        self.cluster_vocab = pd.DataFrame(list(vocab.items()), columns=["id", "name"])

        with open("povinn.pickle", "wb") as f:
            pickle.dump(self.povinn, f, protocol=pickle.HIGHEST_PROTOCOL)
        with open("val.pickle", "wb") as f:
            pickle.dump(self.val, f, protocol=pickle.HIGHEST_PROTOCOL)
        with open("user.pickle", "wb") as f:
            pickle.dump(self.user, f, protocol=pickle.HIGHEST_PROTOCOL)
        with open("train.pickle", "wb") as f:
            pickle.dump(self.train, f, protocol=pickle.HIGHEST_PROTOCOL)
        with open("finished.pickle", "wb") as f:
            pickle.dump(self.finished, f, protocol=pickle.HIGHEST_PROTOCOL)
        with open("cluster_vocab.pickle", "wb") as f:
            pickle.dump(self.cluster_vocab, f, protocol=pickle.HIGHEST_PROTOCOL)

    def filter_out_finished(self, pred, finished):
        return [cid for cid in pred if cid not in finished]

    def group_by_cluster(self, courses, clusters, k=None):
        return group_by_cluster(courses, clusters, k)

    def get_user_soident(self, user: User):
        soident = user.id
        if user.fetch:
            if user.id.lower() == "random":
                uid = (
                    self.val["user_id"]
                    .drop_duplicates()
                    .sample(1, random_state=RND_STATE + self.counter)
                    .iloc[0]
                )
                soident = self.user[self.user["user_id"] == uid]["soident"].iloc[0]
                self.counter += 1
        else:
            soident = ""
        return soident

    def get_user_info(self, user: User, soident: str):
        type = ""
        sobor = ""
        degree_plan = user.degree_plan
        if user.fetch:
            df = self.user[self.user["soident"] == int(soident)]
            type = df["sdruh"].iloc[0]
            sobor = df["sobor_nazev"].iloc[0]
            degree_plan = df["splan"].iloc[0]
        return type, sobor, degree_plan

    def get_year_finished(self, user: User, soident: str):
        year = -1
        finished = user.finished()
        if user.fetch:
            user_id = self.user[self.user["soident"] == int(soident)]["user_id"].iloc[0]
            finished_df = self.train[self.train["user_id"] == user_id]
            finished_df = finished_df.merge(self.povinn, on="course_id")
            year = finished_df["zroc"].max()
            finished = finished_df["povinn"].to_list()

        return year, finished

    def get_expected(self, user: User, soident: str):
        expected = []
        if user.fetch:
            user_id = self.user[self.user["soident"] == int(soident)]["user_id"].iloc[0]
            expected_df = self.val[self.val["user_id"] == user_id]
            expected_df = expected_df.merge(self.povinn, on="course_id")
            expected = expected_df["povinn"].to_list()

        return expected

    def get_degree_plan_courses(self, degree_plan: str):
        dp = self.data.stud_plan
        dp = dp[dp["plan_code"] == degree_plan]
        dp = dp["code"].unique().tolist()
        return dp

    def get_user_id(self, user: User, soident: str):
        uid = None
        if user.fetch:
            uid = self.user[self.user["soident"] == int(soident)]["user_id"].iloc[0]
        return uid

    def user_interaction_povinn(self):
        def sql_builder(with_expr, conn):
            def sql_executor(table):
                df = pd.read_sql(with_expr + f"SELECT * FROM {table}", conn)
                return df

            return sql_executor

        with_expr = """
            WITH istudium AS (
                SELECT
                    soident, sident, sdruh, srokp, sobor, o.nazev sobor_nazev, splan
                FROM studium s
                LEFT JOIN obor o ON s.sobor = o.kod
                WHERE s.sobor like 'I%'
                AND sstav NOT IN ('Z', 'U')
                ORDER BY soident, sident, sdruh, srokp, sobor, o.nazev
            ),
            tmp_interactions AS (
                SELECT
                    soident, sident, zpovinn, zskr::INT, zroc, zsem
                FROM istudium
                LEFT JOIN zkous z ON istudium.sident = z.zident
                WHERE z.zsplcelk = 'S'
                ORDER BY soident, sident, zpovinn, zskr, zroc, zsem
            ),
            interactions AS (
                SELECT DISTINCT
                    i1.soident, i1.sident, i2.zpovinn povinn, i2.zskr, i2.zroc, i2.zsem
                FROM tmp_interactions i1
                LEFT JOIN tmp_interactions i2 ON i1.soident = i2.soident
                ORDER BY i1.soident, i1.sident, i2.zpovinn, i2.zskr, i2.zroc, i2.zsem
            ),
            povinn AS (
                SELECT DISTINCT
                    p.povinn, p.pnazev, panazev, p.pgarant
                FROM interactions i
                LEFT JOIN povinn p ON i.povinn = p.povinn
                ORDER BY p.povinn
            )
        """
        conn = psycopg2.connect(
            dbname=os.getenv("POSTGRES_DB", "recsis"),
            user="recommender",
            host=os.getenv("POSTGRES_HOST", "localhost"),
            password=os.environ["RECSIS_RECOMMENDER_DB_PASS"],
            port=5432,
        )
        load_df = sql_builder(with_expr, conn)
        user = load_df("istudium")
        interactions = load_df("interactions")
        povinn = load_df("povinn")
        conn.close()
        return user, interactions, povinn

    def dataset(self):
        user, interaction, povinn = self.user_interaction_povinn()

        user = user.reset_index().rename(columns={"index": "user_id"})
        # user["sobor_embed"] = list(sbert_embed(user["sobor_nazev"]))
        povinn = povinn.reset_index().rename(columns={"index": "course_id"})
        # povinn["pnazev_embed"] = list(sbert_embed(povinn["pnazev"]))

        interaction = interaction.merge(user[["sident", "user_id"]], on="sident")
        interaction = interaction.merge(povinn[["povinn", "course_id"]], on="povinn")
        interaction = interaction[["user_id", "course_id", "zskr", "zroc"]]

        return user, interaction, povinn

    def split(self, interaction, val_ratio, split_year=2024):
        # Train data are all interactions before split_year
        train = interaction[interaction["zskr"] < split_year]

        # Test data are all interactions after split_year (including split_year)
        year_bitmap = interaction["zskr"] >= split_year

        # Split test data using val_ratio into validation and test sets by user_id randomly
        test_user_id = interaction[year_bitmap]["user_id"].drop_duplicates()
        val_user_id = test_user_id.sample(frac=val_ratio, random_state=RND_STATE)
        val_user_bitmap = interaction["user_id"].isin(val_user_id)
        test_user_bitmap = interaction["user_id"].isin(
            test_user_id.drop(val_user_id.index)
        )

        val = interaction[year_bitmap & val_user_bitmap]
        test = interaction[year_bitmap & test_user_bitmap]
        # val = interaction[test_bitmap].sample(frac=val_ratio)
        # test = interaction[test_bitmap].drop(val.index)

        return train, val, test


def user_interaction_povinn():
    def sql_builder(with_expr, conn):
        def sql_executor(table):
            df = pd.read_sql(with_expr + f"SELECT * FROM {table}", conn)
            return df

        return sql_executor

    with_expr = """
        WITH istudium AS (
            SELECT
                soident, sident, sdruh, srokp, sobor, o.nazev sobor_nazev, splan
            FROM studium s
            LEFT JOIN obor o ON s.sobor = o.kod
            WHERE s.sobor like 'I%'
            AND sstav NOT IN ('Z', 'U')
            ORDER BY soident, sident, sdruh, srokp, sobor, o.nazev
        ),
        tmp_interactions AS (
            SELECT
                soident, sident, zpovinn, zskr::INT, zroc, zsem
            FROM istudium
            LEFT JOIN zkous z ON istudium.sident = z.zident
            WHERE z.zsplcelk = 'S'
            ORDER BY soident, sident, zpovinn, zskr, zroc, zsem
        ),
        interactions AS (
            SELECT DISTINCT
                i1.soident, i1.sident, i2.zpovinn povinn, i2.zskr, i2.zroc, i2.zsem
            FROM tmp_interactions i1
            LEFT JOIN tmp_interactions i2 ON i1.soident = i2.soident
            ORDER BY i1.soident, i1.sident, i2.zpovinn, i2.zskr, i2.zroc, i2.zsem
        ),
        povinn AS (
            SELECT DISTINCT
                p.povinn, p.pnazev, panazev, p.pgarant
            FROM interactions i
            LEFT JOIN povinn p ON i.povinn = p.povinn
            ORDER BY p.povinn
        )
    """
    conn = psycopg2.connect(
        dbname=os.getenv("POSTGRES_DB", "recsis"),
        user="recommender",
        host=os.getenv("POSTGRES_HOST", "localhost"),
        password=os.environ["RECSIS_RECOMMENDER_DB_PASS"],
        port=5432,
    )
    load_df = sql_builder(with_expr, conn)
    user = load_df("istudium")
    interactions = load_df("interactions")
    povinn = load_df("povinn")
    conn.close()
    return user, interactions, povinn


def dataset():
    user, interaction, povinn = user_interaction_povinn()

    user = user.reset_index().rename(columns={"index": "user_id"})
    # user["sobor_embed"] = list(sbert_embed(user["sobor_nazev"]))
    povinn = povinn.reset_index().rename(columns={"index": "course_id"})
    # povinn["pnazev_embed"] = list(sbert_embed(povinn["pnazev"]))

    interaction = interaction.merge(user[["sident", "user_id"]], on="sident")
    interaction = interaction.merge(povinn[["povinn", "course_id"]], on="povinn")
    interaction = interaction[["user_id", "course_id", "zskr", "zroc"]]

    return user, interaction, povinn


if __name__ == "__main__":
    main()
