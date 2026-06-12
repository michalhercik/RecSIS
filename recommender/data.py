import os
import pickle
from pathlib import Path

import numpy as np
import pandas as pd
import psycopg2
from evaluation.embedder import sbert_embed


class TrainData:
    user: pd.DataFrame
    finished: pd.DataFrame
    povinn: pd.DataFrame
    train: pd.DataFrame
    stud_plan: pd.DataFrame
    val: pd.DataFrame
    test: pd.DataFrame
    trida: pd.DataFrame
    klas: pd.DataFrame
    ucit: pd.DataFrame
    categories: pd.DataFrame
    no_history_user_embed: np.ndarray

    rand_soident_counter: int = 0

    CACHE_FILENAME = "dataset_cache.pkl"
    VAL_RATIO = 0.2

    def __init__(self, rnd_state: int):
        self.rnd_state = rnd_state

    def fit(self, cache_root="train-data"):
        self.cache_root = Path(cache_root)
        self.cache_root.mkdir(parents=True, exist_ok=True)

        self.cache_path = self.cache_root / self.CACHE_FILENAME

        if self.cache_path.exists():
            with open(self.cache_path, "rb") as f:
                data = pickle.load(f)
            self.__dict__.update(data)
            return

        user, finished, povinn, stud_plan, klas, trida, ucit, categories = (
            self.dataset()
        )
        train, val, test = self.split(finished, self.VAL_RATIO, 2024)

        user = user.merge(
            train.merge(povinn[["course_id", "embed"]], on="course_id")
            .groupby("user_id")
            .agg(
                {
                    "user_id": "first",
                    "course_id": set,
                    "embed": lambda x: np.mean(x.values, axis=0),
                }
            )[["user_id", "embed"]],
            left_on="user_id",
            right_index=True,
            how="left",
        )[
            [
                "user_id",
                "soident",
                "sident",
                "sdruh",
                "srokp",
                "sobor",
                "sobor_nazev",
                "splan",
                "embed",
            ]
        ]
        no_history_user_embed = list(sbert_embed([""]))[0]

        def _fill_embed(x):
            na = pd.isna(x)
            if isinstance(na, (np.ndarray, list, pd.Series)):
                na = bool(np.all(na))
            if na:
                return no_history_user_embed
            return x

        user["embed"] = user["embed"].apply(_fill_embed)

        data = {
            "user": user,
            "povinn": povinn,
            "val": val,
            "test": test,
            "train": train,
            "finished": finished,
            "stud_plan": stud_plan,
            "klas": klas,
            "trida": trida,
            "ucit": ucit,
            "categories": categories,
            "no_history_user_embed": no_history_user_embed,
        }

        self.__dict__.update(data)
        with open(self.cache_path, "wb") as f:
            pickle.dump(data, f, protocol=pickle.HIGHEST_PROTOCOL)

    def get_finished(self, soident: str) -> list[str]:
        finished = self.finished_df(soident)["povinn"].to_list()
        return finished

    def finished_df(self, soident: str) -> pd.DataFrame:
        user_id = self.user[self.user["soident"] == int(soident)]["user_id"].iloc[0]
        finished_df = self.train[self.train["user_id"] == user_id]
        finished_df = finished_df.merge(self.povinn, on="course_id")
        return finished_df.drop(columns=["user_id", "zskr", "zroc"])

    def degree_plan_courses_by_soident(self, soident: str) -> list[str]:
        degree_plan = self.user[self.user["soident"] == int(soident)]["splan"].iloc[0]
        dp = self.degree_plan_courses_by_code(degree_plan)
        return dp

    def degree_plan_courses_by_code(self, plan_code: str) -> list[str]:
        dp = self.stud_plan[self.stud_plan["plan_code"] == plan_code]
        dp = dp["code"].unique().tolist()
        return dp

    def get_expected(self, soident: str) -> list[str]:
        user_id = self.user[self.user["soident"] == int(soident)]["user_id"]
        if user_id.empty:
            return []
        user_id = user_id.iloc[0]
        expected_df = self.val[self.val["user_id"] == user_id]
        expected_df = expected_df.merge(self.povinn, on="course_id")
        expected = expected_df["povinn"].to_list()
        return expected

    def get_type(self, soident: str) -> str:
        df = self.user[self.user["soident"] == int(soident)]
        return df["sdruh"].iloc[0]

    def get_sobor(self, soident: str) -> str:
        df = self.user[self.user["soident"] == int(soident)]
        return df["sobor_nazev"].iloc[0]

    def get_degree_plan(self, soident: str) -> str:
        df = self.user[self.user["soident"] == int(soident)]
        return df["splan"].iloc[0]

    def get_year_finished(self, soident: str):
        user_id = self.user[self.user["soident"] == int(soident)]["user_id"].iloc[0]
        finished_df = self.val[self.val["user_id"] == user_id]
        year = finished_df["zroc"].max()
        return year

    def get_user(self, soident: str) -> pd.DataFrame:
        return self.user[self.user["soident"] == soident]

    def rand_soident_from_dev(self) -> str:
        uid = (
            self.val["user_id"]
            .drop_duplicates()
            .sample(1, random_state=self.rnd_state + self.rand_soident_counter)
            .iloc[0]
        )
        soident = self.user[self.user["user_id"] == uid]["soident"].iloc[0]
        self.rand_soident_counter += 1
        return soident

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
                    p.povinn, p.pnazev, panazev, p.pgarant, vucit1, vucit2, vucit3
                FROM interactions i
                LEFT JOIN povinn p ON i.povinn = p.povinn
                WHERE p.pgarant != '32-STUD'
                ORDER BY p.povinn
            ),
            pamela AS (
                SELECT pamela.povinn, pamela.typ, pamela.memo
                FROM povinn
                LEFT JOIN pamela ON pamela.povinn = povinn.povinn
                where pamela.jazyk = 'ENG' and pamela.typ in ('A', 'S')
            ),
            filtered_klas AS (
                SELECT povinn, kod, nazev FROM klas
                --WHERE nazev NOT IN ('Předměty obecného základu')
            ),
            filtered_trida AS (
                SELECT * FROM trida
                WHERE nazev NOT LIKE 'M Bc.%'
                AND nazev NOT LIKE 'M Mgr.%'
                --AND nazev NOT IN ('Informatika Bc.', 'Informatika Mgr. - volitelný', 'volitelný', 'Všeobecné')
            ),
            ucit AS (
                SELECT * FROM (
                    SELECT povinn, vucit1 kod, vucit1 nazev FROM povinn
                    UNION
                    SELECT povinn, vucit2 kod, vucit2 nazev FROM povinn
                    UNION
                    SELECT povinn, vucit3 kod, vucit3 nazev FROM povinn
                ) WHERE kod IS NOT NULL
            ),
            categories AS (
                SELECT * FROM filtered_klas
                UNION
                SELECT * FROM filtered_trida
                UNION
                SELECT povinn, pgarant kod, pgarant nazev FROM povinn
                UNION
                SELECT * FROM ucit
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
        stud_plan = sql_builder("", conn)("stud_plan")
        klas = load_df("klas")
        trida = load_df("filtered_trida")
        ucit = load_df("ucit")
        categories = load_df("categories")
        pamela = load_df("pamela")
        conn.close()
        return (
            user,
            interactions,
            povinn,
            stud_plan,
            klas,
            trida,
            ucit,
            categories,
            pamela,
        )

    def dataset(self):
        user, interaction, povinn, stud_plan, klas, trida, ucit, categories, pamela = (
            self.user_interaction_povinn()
        )

        user = user.reset_index().rename(columns={"index": "user_id"})
        povinn = povinn.reset_index().rename(columns={"index": "course_id"})

        pamela = pamela.pivot_table(
            index="povinn", columns="typ", values="memo", aggfunc="first"
        )
        povinn = povinn.merge(pamela, on="povinn", how="left")
        embed_src = povinn.apply(
            lambda x: f"{x['panazev']}: {x['A']}\n{x['S']}", axis=1
        )
        povinn["embed"] = list(sbert_embed(embed_src))

        interaction = interaction.merge(user[["sident", "user_id"]], on="sident")
        interaction = interaction.merge(povinn[["povinn", "course_id"]], on="povinn")
        interaction = interaction[["user_id", "course_id", "zskr", "zroc"]]

        return user, interaction, povinn, stud_plan, klas, trida, ucit, categories

    def split(self, interaction, val_ratio, split_year=2024):
        # Train data are all interactions before split_year
        train = interaction[interaction["zskr"] < split_year]

        # Test data are all interactions after split_year (including split_year)
        year_bitmap = interaction["zskr"] >= split_year

        # Split test data using val_ratio into validation and test sets by user_id randomly
        test_user_id = interaction[year_bitmap]["user_id"].drop_duplicates()
        val_user_id = test_user_id.sample(frac=val_ratio, random_state=self.rnd_state)
        val_user_bitmap = interaction["user_id"].isin(val_user_id)
        test_user_bitmap = interaction["user_id"].isin(
            test_user_id.drop(val_user_id.index)
        )

        val = interaction[year_bitmap & val_user_bitmap]
        test = interaction[year_bitmap & test_user_bitmap]
        # val = interaction[test_bitmap].sample(frac=val_ratio)
        # test = interaction[test_bitmap].drop(val.index)

        return train, val, test
