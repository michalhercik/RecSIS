import os

import pandas as pd
import psycopg2

from algo.base import Algorithm

RND_STATE = 42
VAL_RATIO = 0.2


class TrainData(Algorithm):
    def user_interaction_povinn(self):
        def sql_builder(with_expr, conn):
            def sql_executor(table):
                df = pd.read_sql(with_expr + f"SELECT * FROM {table}", conn)
                return df

            return sql_executor

        with_expr = """
            WITH istudium AS (
                SELECT
                    soident, sident, sdruh, srokp, sobor, o.nazev sobor_nazev
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
            host="postgres",
            password=os.environ["RECSIS_RECOMMENDER_DB_PASS"],
            port=5432,
        )
        load_df = sql_builder(with_expr, conn)
        user = load_df("istudium")
        interactions = load_df("interactions")
        povinn = load_df("povinn")
        conn.close()
        return user, interactions, povinn

    def dataset(self, force_load=True):
        user, interaction, povinn = self.user_interaction_povinn()

        user = user.reset_index().rename(columns={"index": "user_id"})
        # user["sobor_embed"] = list(sbert_embed(user["sobor_nazev"]))
        povinn = povinn.reset_index().rename(columns={"index": "course_id"})
        # povinn["pnazev_embed"] = list(sbert_embed(povinn["pnazev"]))

        interaction = interaction.merge(user[["sident", "user_id"]], on="sident")
        interaction = interaction.merge(povinn[["povinn", "course_id"]], on="povinn")
        interaction = interaction[["user_id", "course_id", "zskr"]]

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
