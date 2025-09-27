import pandas as pd
from sqlalchemy import create_engine


class DataRepository:
    def __init__(self):
        pass

    def get_degree_plan(self, code: str, year: int):
        engine = create_engine("postgresql://recommender:recommender@localhost:5432/recsis")
        with engine.connect() as conn:
            dp = pd.read_sql_query(
                f"SELECT * FROM webapp.degree_plans WHERE plan_code = '{code}' AND plan_year = {year}", 
                conn
            )
        return dp

    def get_courses(self):
        engine = create_engine("postgresql://recommender:recommender@localhost:5432/recsis")
        with engine.connect() as conn:
            courses = pd.read_sql_query("SELECT * FROM webapp.courses", conn)
        return courses