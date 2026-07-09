import os
import threading
from functools import partial
from multiprocessing.sharedctypes import Synchronized

import pandas as pd
from sqlalchemy import create_engine


class DataReader:
    def __init__(self):
        host = os.environ.get("POSTGRES_HOST", "localhost")
        self.engine = create_engine(
            f"postgresql://recommender:recommender@{host}:5432/recsis"
        )

    def fetch_table(self, table_name: str):
        with self.engine.connect() as conn:
            return pd.read_sql_table(table_name, conn)

    def fetch_query(self, query: str):
        with self.engine.connect() as conn:
            return pd.read_sql_query(query, conn)


class LazyLoader:
    """Descriptor for lazy-loading tables into each DataRepository instance."""

    def __init__(self, table_name: str, reader: "DataReader"):
        self.table_name = table_name
        self.reader = reader
        self._lock = threading.Lock()

    def __get__(self, instance, owner):
        if instance is None:
            return self

        attr_name = f"__cached_{self.table_name}"

        # already loaded?
        if hasattr(instance, attr_name):
            return getattr(instance, attr_name)

        # load thread-safe
        with self._lock:
            if not hasattr(instance, attr_name):
                table = self.reader.fetch_table(self.table_name)
                setattr(instance, attr_name, table)

        return getattr(instance, attr_name)

    class LazyQueryLoader:
        """Descriptor for lazy-loading tables into each DataRepository instance."""

        def __init__(self, query: str, reader: "DataReader"):
            self.query = query
            self.reader = reader
            self._lock = threading.Lock()

        def __get__(self, instance, owner):
            if instance is None:
                return self

            attr_name = f"__cached_{self.query}"

            # already loaded?
            if hasattr(instance, attr_name):
                return getattr(instance, attr_name)

            # load thread-safe
            with self._lock:
                if not hasattr(instance, attr_name):
                    table = self.reader.fetch_query(self.query)
                    setattr(instance, attr_name, table)

            return getattr(instance, attr_name)


class DataRepository:
    """Repository with lazily loaded tables."""

    _reader = DataReader()

    povinn = LazyLoader("povinn", _reader)
    zkous = LazyLoader("zkous", _reader)
    studium = LazyLoader("studium", _reader)
    stud_plan = LazyLoader("stud_plan", _reader)
    pamela = LazyLoader("pamela", _reader)


def user_interaction_povinn():
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
        host="localhost",
        password=os.environ["RECSIS_RECOMMENDER_DB_PASS"],
        port=5432,
    )
    load_df = sql_builder(with_expr, conn)
    user = load_df("istudium")
    interactions = load_df("interactions")
    povinn = load_df("povinn")
    conn.close()
    return user, interactions, povinn
