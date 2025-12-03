import os

import pandas as pd
from sqlalchemy import create_engine


class DataRepository:
    def __init__(self):
        host = os.environ.get("POSTGRES_HOST", "localhost")
        self.__engine = create_engine(
            f"postgresql://recommender:recommender@{host}:5432/recsis"
        )
        self.povinn = self.__fetch_table("povinn")
        self.zkous = self.__fetch_table("zkous")
        self.studium = self.__fetch_table("studium")
        self.stud_plan = self.__fetch_table("stud_plan")
        self.pamela = self.__fetch_table("pamela")

    def sql_query(self, query: str, params=None):
        with self.__engine.connect() as conn:
            result = pd.read_sql_query(query, params, conn)
        return result

    def __fetch_table(self, table_name: str):
        with self.__engine.connect() as conn:
            table = pd.read_sql_table(table_name, conn)
        return table
