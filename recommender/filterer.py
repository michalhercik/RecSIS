from data import TrainData
from user import User


class Filterer:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def filter(self, user: User, courses: list[str]) -> list[str]:
        raise NotImplementedError()


class FinishedFilter(Filterer):
    def filter(self, user: User, courses: list[str]) -> list[str]:
        """
        Filters out courses that the user has already finished.
        param:
            courses: list of course titles to filter
        return:
            list of course titles that the user has not finished
        """
        finished = user.finished()
        if user.fetch:
            finished = self.train_data.get_finished(user.id)
        return [cid for cid in courses if cid not in finished]
