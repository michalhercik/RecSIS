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


class DataRepository:
    """Repository with lazily loaded tables."""

    _reader = DataReader()

    povinn = LazyLoader("povinn", _reader)
    zkous = LazyLoader("zkous", _reader)
    studium = LazyLoader("studium", _reader)
    stud_plan = LazyLoader("stud_plan", _reader)
    pamela = LazyLoader("pamela", _reader)
