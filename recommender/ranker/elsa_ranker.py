import typing

import numpy as np
import pandas as pd
import torch
from data import TrainData
from elsa import ELSA
from ranker.ranker import Ranker
from user import User


class Elsa(Ranker):
    def __init__(self, train_data: TrainData):
        super().__init__(train_data)
        self.set_train_params(
            factors=16,
            num_epochs=50,
            learning_rate=1e-2,
        )

    def fit(self) -> None:
        train_im = self.interaction_matrix(self.train_data.user, self.train_data.train)

        self.model = ELSA(
            n_items=self.train_data.povinn.shape[0],
            device=self.device,
            n_dims=self.factors,
            lr=self.learning_rate,
        )
        self.model.fit(
            train_im.values,
            batch_size=self.train_data.train.shape[0],
            epochs=self.num_epochs,
            shuffle=False,
            verbose=False,
        )

    def rank(self, user: User) -> list[str]:
        bp_im = self.__interaction_matrix_from(user)
        pred = self.model.predict(bp_im.values, batch_size=1)
        topk = torch.topk(pred, k=pred.shape[1], sorted=True)
        pred = self.train_data.povinn["povinn"].iloc[topk.indices[0]].to_list()
        return pred

    def similar_items(
        self,
        N: int,
        batch_size: int,
        sources: typing.Union[np.ndarray, torch.Tensor] = None,
        candidates: typing.Union[np.ndarray, torch.Tensor] = None,
        verbose: bool = True,
    ) -> tuple:
        return self.model.similar_items(
            N=N,
            batch_size=batch_size,
            sources=sources,
            candidates=candidates,
            verbose=verbose,
        )

    def set_train_params(
        self, factors, num_epochs, learning_rate, device=torch.device("cpu")
    ):
        self.factors = factors
        self.num_epochs = num_epochs
        self.learning_rate = learning_rate
        self.device = device

    def __interaction_matrix_from(self, user: User):
        bp_im = None
        user_id = None
        finished = None
        if user.fetch:
            user_id = self.train_data.user[
                self.train_data.user["soident"] == int(user.id)
            ]
            finished = self.train_data.train[
                self.train_data.train["user_id"] == user_id["user_id"].iloc[0]
            ]
            finished = finished.merge(self.train_data.povinn, on="course_id")
            # bp_im = self.interaction_matrix(user_id, finished, self.povinn)
        else:
            user_id = pd.DataFrame({"user_id": [user.id]})
            finished = user.blueprint_to_df()
            finished = finished.merge(
                self.train_data.povinn, left_on="course", right_on="povinn"
            )
            finished["user_id"] = user.id
        bp_im = self.interaction_matrix(user_id, finished)
        return bp_im

    def interaction_matrix(self, user, finished):
        im = pd.crosstab(finished["user_id"], finished["course_id"])
        im = im.reindex(
            index=user["user_id"],
            columns=self.train_data.povinn["course_id"],
            fill_value=0,
        )
        return im
