import pandas as pd
from sqlalchemy import create_engine

class DataRepository:
    def __init__(self):
        self.__engine = create_engine("postgresql://recommender:recommender@localhost:5432/recsis")
        self.povinn = self.__fetch_table("povinn")
        self.zkous = self.__fetch_table("zkous")
        self.studium = self.__fetch_table("studium")
        self.stud_plan = self.__fetch_table("stud_plan")
    
    def sql_query(self, query: str):
        with self.__engine.connect() as conn:
            result = pd.read_sql_query(query, conn)
        return result

    def __fetch_table(self, table_name: str):
        with self.__engine.connect() as conn:
            table = pd.read_sql_table(table_name, conn)
        return table