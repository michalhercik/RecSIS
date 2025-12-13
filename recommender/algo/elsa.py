import numpy as np
import pandas as pd
import torch
from algo.base import Algorithm
from elsa import ELSA
from user import User


class Elsa(Algorithm):
    def fit(self):
        device = torch.device("cuda")
        # TODO: implement X
        raise NotImplementedError("Please implement this method")
        X = []
        factors = 256
        num_epochs = 5
        batch_size = 128
        self.fit_with_params(X, factors, num_epochs, batch_size, device)

    def recommend(self, user: User, limit: int) -> list[str]:
        # TODO: implement
        # transform user
        raise NotImplementedError("Please implement this method")

    def predict(self, X, batch_size):
        predictions = self.model.predict(X, batch_size=batch_size)
        # predictions = ((X @ self.A) @ (self.A.T)) - X
        return predictions

    def similar_items(self, itemids):
        related = self.model.similar_items(N=100, batch_size=128, sources=itemids)
        return related

    def fit_with_params(
        self,
        X,
        validation_data,
        factors,
        num_epochs,
        batch_size,
        shuffle,
        device,
        lr=0.1,
    ):
        items_cnt = X.shape[1]
        self.model = ELSA(n_items=items_cnt, device=device, n_dims=factors, lr=lr)
        self.model.fit(
            X,
            validation_data=validation_data,
            batch_size=batch_size,
            epochs=num_epochs,
            shuffle=shuffle,
        )
        # self.A = (
        #     torch.nn.functional.normalize(self.model.get_items_embeddings(), dim=-1)
        #     .cpu()
        #     .numpy()
        # )
