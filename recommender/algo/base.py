import pandas as pd
from data_repository import DataRepository
from user import User

class Algorithm:
    def __init__(self, data: DataRepository):
        self.data = data

    def recommend(self, user: User, limit: int) -> list[str]:
        raise NotImplementedError("Please Implement this method")

    def fit(self):
        raise NotImplementedError("Algorithm.fit() is not implemented")
