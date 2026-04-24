from data import TrainData
from user import User


class Ranker:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def fit(self) -> None:
        raise NotImplementedError()

    def rank(self, user: User) -> list[str]:
        raise NotImplementedError()
