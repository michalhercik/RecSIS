import os
import pickle

from data import TrainData
from user import User


class Ranker:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def fit(self) -> None:
        raise NotImplementedError()

    def rank(self, user: User) -> list[str]:
        raise NotImplementedError()


def cached(retrieve, file, condition=None):
    if condition is None or condition():
        if os.path.exists(file):
            with open(file, "rb") as f:
                return pickle.load(f)
    data = retrieve()
    with open(file, "wb") as f:
        pickle.dump(data, f, protocol=pickle.HIGHEST_PROTOCOL)
    return data
