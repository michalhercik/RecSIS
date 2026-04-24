from data import TrainData
from user import User


class Masker:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def fit(self, user: User):
        raise NotImplementedError()

    def mask(self, courses: list[str]) -> list[bool]:
        raise NotImplementedError()


class TruePositivesMasker(Masker):
    def fit(self, user: User):
        self.expected = []
        if user.fetch:
            self.expected = set(self.train_data.get_expected(user.id))

    def mask(self, courses: list[str]) -> list[bool]:
        mask = [course in self.expected for course in courses]
        return mask


class DegreePlanMasker(Masker):
    def fit(self, user: User):
        dp = []
        if user.fetch:
            dp = self.train_data.degree_plan_courses_by_soident(user.id)
        elif user.degree_plan is not None:
            dp = self.train_data.degree_plan_courses_by_code(user.degree_plan)
        self.dp_courses = set(dp)

    def mask(self, courses: list[str]) -> list[bool]:
        mask = [course in self.dp_courses for course in courses]
        return mask


class FalseNegativeMasker(Masker):
    def fit(self, user: User):
        self.expected = []
        if user.fetch:
            self.expected = set(self.train_data.get_expected(user.id))

    def mask(self, courses: list[str]) -> list[bool]:
        mask = [course not in courses for course in self.expected]
        return mask
