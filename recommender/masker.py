from data import TrainData
from user import User


class FittableMasker:
    def __init__(self, train_data: TrainData):
        self.train_data = train_data

    def fit(self, user: User):
        raise NotImplementedError()

    def mask(self, courses: list[str]) -> list[bool]:
        raise NotImplementedError()


class Masker:
    def mask(self, reference: list[str], values: list[str]) -> list[bool]:
        """
        keep order of values
        """
        raise NotImplementedError()


class TruePositivesMasker(FittableMasker):
    def fit(self, user: User):
        self.expected = []
        if user.fetch:
            self.expected = set(self.train_data.get_expected(user.id))

    def mask(self, courses: list[str]) -> list[bool]:
        mask = [course in self.expected for course in courses]
        return mask


class DegreePlanMasker(FittableMasker):
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


class InMasker(Masker):
    def mask(self, reference: list[str], values: list[str]) -> list[bool]:
        mask = [course in reference for course in values]
        return mask


class NotInMasker(InMasker):
    def mask(self, reference: list[str], values: list[str]) -> list[bool]:
        mask = [not m for m in super().mask(reference, values)]
        return mask
