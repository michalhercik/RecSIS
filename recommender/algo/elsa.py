import pandas as pd
import torch
from algo.base import Recommendation, Result
from algo.train import TrainData
from data_repository import DataRepository
from elsa import ELSA
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
        # user, finished, povinn = self.dataset()
        # train, val, test = self.split(finished, VAL_RATIO, 2024)

        super().fit()

        train_im = self.interaction_matrix(self.user, self.train, self.povinn)

        self.model = ELSA(
            n_items=self.povinn.shape[0],
            device=self.device,
            n_dims=self.factors,
            lr=self.learning_rate,
        )
        self.model.fit(
            train_im.values,
            batch_size=self.train.shape[0],
            epochs=self.num_epochs,
            shuffle=False,
        )
        # self.user = user
        # self.finished_train = train
        # self.finished_val = val
        # self.povinn = povinn

    # 909953

    def recommend(self, user: User, limit: int) -> Result:
        result = Result()
        result.soident = self.get_user_soident(user)
        result.type, result.sobor, result.degree_plan = self.get_user_info(
            user, result.soident
        )
        result.year_of_study, result.finished = self.get_year_finished(
            user, result.soident
        )
        result.expected = self.get_expected(user, result.soident)

        # if user.fetch:
        #     bp_im, finished, dp_code, expected, user_id, sobor, sdruh, zroc = self.interaction_matrix_from_train_data(user)
        # else:
        #     bp_im, finished, dp_code  = self.interaction_matrix_from_user(user)
        #     user_id = user.id
        #     expected = list()
        bp_im = self.interaction_matrix_from(user, result.soident)
        pred = self.model.predict(bp_im.values, batch_size=1)
        topk = torch.topk(pred, k=pred.shape[1], sorted=True)
        pred = self.povinn["povinn"].iloc[topk.indices[0]].to_list()

        # pred = [i for i in predictions if i not in set(finished)]
        pred = self.filter_out_finished(pred, result.finished)
        pred = self.povinn[self.povinn["povinn"].isin(pred)]
        pred = self.group_by_cluster(pred["povinn"], pred["cluster"], k=limit)

        # dp = self.data.stud_plan
        # dp = dp[dp["plan_code"] == dp_code]
        # dp = dp["code"].unique().tolist()
        # predictions = [i for i in predictions if i not in set(dp)]

        # result.recommended = pred[:limit]
        result.recommended = pred
        dp_courses = self.get_degree_plan_courses(result.degree_plan)
        result.generate_masks(dp_courses)
        # target = [True if i in set(expected) else False for i in predictions]
        # result = Result(
        #     soident=user_id,
        #     sobor=sobor,
        #     degree_plan=dp_code,
        #     year_of_study=zroc,
        #     type="Bachelor" if sdruh == "B" else "Master",
        #     finished=finished,
        #     finished_in_degree_plan=[True if i in set(dp) else False for i in finished],
        #     recommended_in_degree_plan=[True if i in set(dp) else False for i in predictions],
        #     expected_in_degree_plan=[True if i in set(dp) else False for i in expected]
        # ).recommended_and_expected(
        #         recommended=predictions,
        #         expected=expected
        #     )
        return result
        # return Recommendation(predictions, target, finished, expected)

    def interaction_matrix(self, user, finished, povinn):
        im = pd.crosstab(finished["user_id"], finished["course_id"])
        im = im.reindex(
            index=user["user_id"], columns=povinn["course_id"], fill_value=0
        )
        return im

    def interaction_matrix_from(self, user: User, soident: str):
        bp_im = None
        user_id = None
        finished = None
        if user.fetch:
            user_id = self.user[self.user["soident"] == int(soident)]
            finished = self.train[self.train["user_id"] == user_id["user_id"].iloc[0]]
            finished = finished.merge(self.povinn, on="course_id")
            # bp_im = self.interaction_matrix(user_id, finished, self.povinn)
        else:
            user_id = pd.DataFrame({"user_id": [user.id]})
            finished = user.blueprint_to_df()
            finished = bp.merge(self.povinn, left_on="course", right_on="povinn")
            finished["user_id"] = user.id
        bp_im = self.interaction_matrix(user_id, finished, self.povinn)
        return bp_im
        # bp_im = pd.crosstab(bp["user_id"], bp["course_id"])
        # bp_im = bp_im.reindex(
        #     index=bp_im.index, columns=self.povinn["course_id"], fill_value=0
        # )

    def interaction_matrix_from_train_data(self, user: User):
        if user.id.lower() == "random":
            uid = (
                self.finished_val["user_id"]
                .drop_duplicates()
                .sample(1, random_state=RND_STATE + self.counter)
                .iloc[0]
            )
            self.counter += 1
            user.id = self.user[self.user["user_id"] == uid]["soident"].iloc[0]
        user_id = self.user[self.user["soident"] == int(user.id)]
        finished = self.finished_train[
            self.finished_train["user_id"] == user_id["user_id"].iloc[0]
        ]
        finished = finished.merge(self.povinn, on="course_id")
        expected = self.finished_val[
            self.finished_val["user_id"] == user_id["user_id"].iloc[0]
        ]
        expected = expected.merge(self.povinn, on="course_id")
        bp_im = self.interaction_matrix(user_id, finished, self.povinn)
        return (
            bp_im,
            finished["povinn"].to_list(),
            user_id["splan"].iloc[0],
            expected["povinn"].to_list(),
            user.id,
            user_id["sobor_nazev"].iloc[0],
            user_id["sdruh"].iloc[0],
            finished["zroc"].max(),
        )

    def interaction_matrix_from_user(self, user: User):
        bp = user.blueprint_to_df()
        bp = bp.merge(self.povinn, left_on="course", right_on="povinn")
        bp["user_id"] = user.id
        bp_im = pd.crosstab(bp["user_id"], bp["course_id"])
        bp_im = bp_im.reindex(
            index=bp_im.index, columns=self.povinn["course_id"], fill_value=0
        )
        return bp_im, bp["course"].to_list(), user.degree_plan
