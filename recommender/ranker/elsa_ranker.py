import numpy as np
import pandas as pd
import torch
from data import TrainData
from elsa import ELSA
from ranker.ranker import Ranker
from sklearn.cluster import KMeans
from sklearn.decomposition import PCA
from user import User


class Elsa(Ranker):
    def __init__(self, train_data: TrainData):
        super().__init__(train_data)
        self.set_train_params(
            factors=256, num_epochs=5, learning_rate=1e-2, batch_size=128
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
        )

        self.clusters = self.cluster()

    def cluster(self):
        emb_matrix = self.model.get_items_embeddings(as_numpy=True)
        embeddings = pd.DataFrame(
            {
                "povinn": self.train_data.povinn["povinn"],
                "pnazev": self.train_data.povinn["pnazev"],
                "embedding": list(emb_matrix),
            }
        )

        orig_dim = emb_matrix.shape[1]
        if orig_dim > 50:
            pca = PCA(n_components=50, random_state=42)
            reduced_emb = pca.fit_transform(emb_matrix)
        else:
            reduced_emb = emb_matrix.copy()

        # also keep reduced embeddings in the dataframe for later inspection
        embeddings["reduced_embedding"] = list(reduced_emb)

        # Build a matrix of embeddings (n_items x n_dims)
        # emb_matrix = np.vstack(embeddings["embedding"].values)
        n_items = emb_matrix.shape[0]

        # Heuristic for number of clusters: at least 1, otherwise sqrt(n_items)
        if n_items <= 1:
            n_clusters = 1
        else:
            n_clusters = max(2, int(n_items**0.5))

        # Run k-means clustering on the item embeddings
        if n_clusters == 1:
            labels = np.zeros(n_items, dtype=int)
            cluster_centers = reduced_emb.copy()
        else:
            kmeans = KMeans(n_clusters=n_clusters, random_state=42, n_init=10)
            labels = kmeans.fit_predict(reduced_emb)
            cluster_centers = kmeans.cluster_centers_

        # Attach cluster labels to the embeddings dataframe
        embeddings["cluster"] = labels

        return embeddings

    def rank(self, user: User) -> list[str]:
        bp_im = self.__interaction_matrix_from(user)
        pred = self.model.predict(bp_im.values, batch_size=1)
        topk = torch.topk(pred, k=pred.shape[1], sorted=True)
        pred = self.train_data.povinn["povinn"].iloc[topk.indices[0]].to_list()
        return pred

    def explain(self, courses: list[str]):
        pass

    def set_train_params(
        self, factors, num_epochs, batch_size, learning_rate, device=torch.device("cpu")
    ):
        self.factors = factors
        self.num_epochs = num_epochs
        self.batch_size = batch_size
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


if __name__ == "__main__":
    main()
