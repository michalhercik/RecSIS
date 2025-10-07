import numpy as np
import pandas as pd
from sklearn.neighbors import NearestNeighbors

from algo.base import Algorithm
from user import User

#=========================================================================================
# Example #1: Random
#=========================================================================================

class random(Algorithm):
    def recommend(self, user: User, limit: int) -> list[str]:
        return self.data.povinn.sample(limit)["povinn"].tolist()

#=========================================================================================
# Example #2: kNN
#=========================================================================================
    
class knn(Algorithm):
    def fit(self):
        self.model = NearestNeighbors()
        self.df = self.data.zkous.sample(1000)
        self.df = self.df \
            .pivot_table(
                index="zident", 
                columns="zpovinn", 
                values="zsplcelk", 
                aggfunc= lambda x: "S" if (x == "S").any() else "N"
            ) \
            .fillna('N') \
            .drop_duplicates() \
            .replace({'S': 1, 'N': 0}) 
        self.model.fit(self.df.values, self.df.index)

    def recommend(self, user: User, limit: int) -> list[str]:
        bp_df = user.blueprint_to_df()
        x = self.__blueprint_to_vector(bp_df)
        n_neighbors = min(np.pow(2, limit-1), self.df.shape[0])

        indices = self.model.kneighbors(
            x.reshape(1, -1), 
            n_neighbors=n_neighbors, 
            return_distance=False
        )

        neighbors = self.df.iloc[indices[0]]
        candidates = neighbors.any(axis=0)
        candidates[bp_df["course"]] = False
        candidates = candidates[candidates == True]
        if candidates.shape[0] < limit:
            limit = candidates.shape[0]
        result = candidates.sample(limit)

        return result.index.tolist()
    
    def __blueprint_to_vector(self, blueprint: pd.DataFrame) -> np.ndarray:
        x = pd.Series(0, index=self.df.columns)
        x[blueprint["course"]] = 1
        return x.to_numpy()