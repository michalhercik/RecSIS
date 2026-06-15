import numpy as np
from data import TrainData
from ranker.ranker import Ranker
from user import User


def cos_sim(x1, x2):
    return np.dot(x1, x2) / (np.linalg.norm(x1) * np.linalg.norm(x2))


class UserKNN(Ranker):
    def fit(self) -> None:
        pass

    def rank(self, user: User) -> list[str]:
        if user.fetch:
            embed = self.train_data.user[
                self.train_data.user["soident"] == int(user.id)
            ]["embed"].iloc[0]
        else:
            bp = user.blueprint_to_df()
            bp = bp.merge(self.train_data.povinn, left_on="course", right_on="povinn")
            embed = bp["embed"].mean(axis=0)

        results = self.train_data.user.copy()
        results["sim"] = results.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
        results = results.merge(self.train_data.train, on="user_id")
        results = results.sort_values("sim", ascending=False).drop_duplicates(
            subset="course_id"
        )
        results = results.merge(self.train_data.povinn, on="course_id")
        return results["povinn"].to_list()


class ContentKNN(Ranker):
    def fit(self) -> None:
        pass

    def rank(self, user: User) -> list[str]:
        if user.fetch:
            embed = self.train_data.user[
                self.train_data.user["soident"] == int(user.id)
            ]["embed"].iloc[0]
        else:
            bp = user.blueprint_to_df()
            bp = bp.merge(self.train_data.povinn, left_on="course", right_on="povinn")
            embed = bp["embed"].mean(axis=0)

        results = self.train_data.povinn.copy()
        results["sim"] = results.apply(lambda x: cos_sim(x["embed"], embed), axis=1)
        results = results.sort_values("sim", ascending=False)
        return results["povinn"].to_list()
