from data import TrainData
from user import User


class Explainer:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def fit(self) -> None:
        raise NotImplementedError()

    def explain(self, user: User, courses: list[str]) -> dict[str, str]:
        raise NotImplementedError()


class EmptyExplainer(Explainer):
    def __init__(self):
        pass

    def fit(self) -> None:
        pass

    def explain(self, user: User, courses: list[str]) -> dict[str, str]:
        return {}
