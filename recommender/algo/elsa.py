import pandas as pd
import torch
from elsa import ELSA

from algo.train import TrainData
from algo.base import Recommendation
from data_repository import DataRepository
from user import User

RND_STATE = 42
VAL_RATIO = 0.2


class Elsa(TrainData):
    def __init__(self, data_repository: DataRepository):
        self.set_train_params(
            factors=256, num_epochs=5, learning_rate=1e-2, batch_size=128
        )
        self.counter = 0
        super().__init__(data_repository)

    def set_train_params(
        self, factors, num_epochs, batch_size, learning_rate, device=torch.device("cpu")
    ):
        self.factors = factors
        self.num_epochs = num_epochs
        self.batch_size = batch_size
        self.learning_rate = learning_rate
        self.device = device

    def fit(self):
        user, finished, povinn = self.dataset()
        train, val, test = self.split(finished, VAL_RATIO, 2024)


        train_im = self.interaction_matrix(user, train, povinn)

        self.model = ELSA(
            n_items=povinn.shape[0],
            device=self.device,
            n_dims=self.factors,
            lr=self.learning_rate,
        )
        self.model.fit(
            train_im.values,
            batch_size=train.shape[0],
            epochs=self.num_epochs,
            shuffle=False,
        )
        self.user = user
        self.finished_train = train
        self.finished_val = val
        self.povinn = povinn


    def recommend(self, user: User, limit: int) -> Recommendation:
        if user.fetch:
            bp_im, finished, dp_code, expected = self.interaction_matrix_from_train_data(user)
        else:
            bp_im, finished, dp_code  = self.interaction_matrix_from_user(user)
            expected = list()

        predictions = self.model.predict(bp_im.values, batch_size=1)
        topk = torch.topk(predictions, k=predictions.shape[1], sorted=True)
        predictions = self.povinn["povinn"].iloc[topk.indices[0]].to_list()

        predictions = [i for i in predictions if i not in set(finished)]

        dp = self.data.stud_plan
        dp = dp[dp["plan_code"] == dp_code]
        dp = dp["code"].unique().tolist()
        predictions = [i for i in predictions if i not in set(dp)]
        predictions = predictions[:limit]
        target = [True if i in set(expected) else False for i in predictions]
        return Recommendation(predictions, target, finished, expected)

    def interaction_matrix(self, user, finished, povinn):
        im = pd.crosstab(finished["user_id"], finished["course_id"])
        im = im.reindex(
            index=user["user_id"], columns=povinn["course_id"], fill_value=0
        )
        return im

    def interaction_matrix_from_train_data(self, user: User):
        if user.id.lower() == "random":
            uid = self.finished_val["user_id"].drop_duplicates().sample(1, random_state=RND_STATE + self.counter).iloc[0]
            self.counter += 1
            user.id = self.user[self.user["user_id"] == uid]["soident"].iloc[0]
        user_id = self.user[self.user["soident"] == int(user.id)]
        finished = self.finished_train[self.finished_train["user_id"] == user_id["user_id"].iloc[0]]
        finished = finished.merge(self.povinn, on="course_id")
        expected = self.finished_val[self.finished_val["user_id"] == user_id["user_id"].iloc[0]]
        expected = expected.merge(self.povinn, on="course_id")
        bp_im = self.interaction_matrix(user_id, finished, self.povinn)
        return bp_im, finished["povinn"].to_list(), user.degree_plan, expected["povinn"].to_list()

    def interaction_matrix_from_user(self, user: User):
        bp = user.blueprint_to_df()
        bp = bp.merge(self.povinn, left_on="course", right_on="povinn")
        bp["user_id"] = user.id
        bp_im = pd.crosstab(bp["user_id"], bp["course_id"])
        bp_im = bp_im.reindex(
            index=bp_im.index, columns=self.povinn["course_id"], fill_value=0
        )
        return bp_im, bp["course"].to_list(), user.degree_plan, list()
