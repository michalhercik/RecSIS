from data import TrainData
from ranker.ranker import Ranker
from user import User


class MeiliSearch(Ranker):
    def __init__(self, train_data: TrainData):
        raise NotImplementedError()

    def fit(self) -> None:
        raise NotImplementedError()

    def rank(self, user: User) -> list[str]:
        raise NotImplementedError()
