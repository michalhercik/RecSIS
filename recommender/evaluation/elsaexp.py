import sys

import numpy as np
import pandas as pd
import torch
from scipy.sparse import csr_matrix

sys.path.insert(0, "..")
from algo.elsa import Elsa
from data_repository import DataRepository

DEVICE = "cuda" if torch.cuda.is_available() else "cpu"


class ElsaRelevantExp:
    """
    Results
        - On test data for k=3,10,20,50 is
            - precision: 66,60,52,32 %
            - recall:    10,30,49,69 %
            - avg-prec:  69,68,65,60 %

    Train
                        precision                         recall                         avg_prec
                        mean   std   25%   50%   75%   mean   std   25%   50%   75%     mean   std   25%   50%   75%
    target        k
    next_semester 3     0.20  0.29  0.00  0.00  0.33   0.10  0.18  0.00  0.00  0.17     0.26  0.37  0.00  0.00  0.50
                  10    0.19  0.19  0.00  0.20  0.30   0.32  0.31  0.00  0.29  0.50     0.29  0.30  0.00  0.23  0.49
                  20    0.16  0.13  0.00  0.15  0.25   0.51  0.37  0.00  0.60  0.83     0.28  0.26  0.00  0.24  0.44
                  50    0.09  0.07  0.02  0.08  0.14   0.68  0.41  0.33  0.89  1.00     0.26  0.24  0.03  0.23  0.40
    relevant      3     0.77  0.36  0.67  1.00  1.00   0.15  0.14  0.08  0.12  0.17     0.81  0.36  0.83  1.00  1.00
                  10    0.71  0.34  0.40  0.90  1.00   0.41  0.22  0.27  0.40  0.53     0.79  0.32  0.71  0.98  1.00
                  20    0.61  0.33  0.30  0.70  0.90   0.65  0.25  0.51  0.67  0.82     0.77  0.31  0.65  0.92  0.99
                  50    0.35  0.23  0.16  0.36  0.50   0.86  0.23  0.83  0.93  1.00     0.72  0.30  0.59  0.86  0.94

    Validate
                        precision                         recall                         avg_prec
                        mean   std   25%   50%   75%   mean   std   25%   50%   75%     mean   std   25%   50%   75%
    target        k
    next_semester 3     0.18  0.27  0.00  0.00  0.33   0.09  0.15  0.00  0.00  0.14     0.23  0.35  0.00  0.00  0.33
                  10    0.17  0.19  0.00  0.10  0.30   0.26  0.27  0.00  0.20  0.50     0.26  0.28  0.00  0.20  0.43
                  20    0.13  0.13  0.00  0.10  0.25   0.41  0.35  0.00  0.44  0.71     0.25  0.25  0.00  0.20  0.40
                  50    0.07  0.06  0.00  0.08  0.14   0.58  0.41  0.00  0.73  1.00     0.22  0.22  0.00  0.18  0.36
    relevant      3     0.67  0.38  0.33  0.67  1.00   0.12  0.12  0.06  0.09  0.14     0.72  0.39  0.33  1.00  1.00
                  10    0.61  0.34  0.30  0.70  0.90   0.31  0.18  0.21  0.30  0.41     0.70  0.34  0.47  0.85  0.99
                  20    0.51  0.30  0.20  0.55  0.75   0.50  0.23  0.38  0.51  0.65     0.67  0.32  0.47  0.78  0.92
                  50    0.31  0.21  0.14  0.32  0.42   0.71  0.25  0.64  0.78  0.86     0.61  0.30  0.41  0.71  0.85

    Test
                        precision                        recall                         avg_prec
                        mean   std  25%   50%   75%   mean   std   25%   50%   75%     mean   std   25%   50%   75%
    target        k
    next_semester 3     0.18  0.27  0.00  0.00  0.33   0.08  0.13  0.00  0.00  0.14     0.24  0.35  0.00  0.00  0.50
                  10    0.18  0.18  0.00  0.10  0.30   0.27  0.27  0.00  0.25  0.50     0.27  0.28  0.00  0.20  0.47
                  20    0.15  0.13  0.00  0.15  0.25   0.44  0.35  0.00  0.50  0.75     0.25  0.24  0.00  0.22  0.42
                  50    0.08  0.07  0.00  0.08  0.14   0.60  0.40  0.00  0.75  1.00     0.23  0.22  0.00  0.20  0.39
    relevant      3     0.66  0.40  0.33  0.67  1.00   0.10  0.09  0.05  0.09  0.14     0.69  0.41  0.33  1.00  1.00
                  10    0.60  0.35  0.30  0.70  0.90   0.30  0.19  0.20  0.29  0.40     0.68  0.36  0.46  0.84  1.00
                  20    0.52  0.32  0.20  0.55  0.75   0.49  0.24  0.39  0.52  0.65     0.65  0.34  0.43  0.76  0.95
                  50    0.31  0.21  0.12  0.30  0.42   0.69  0.27  0.61  0.75  0.87     0.60  0.31  0.38  0.69  0.86
    """

    def __init__(self, data: DataRepository, train, validate):
        self.name = "ElsaRelevant"

        X = train[["relevant"]].explode("relevant")
        X = pd.crosstab(X.index, X["relevant"])
        self.columns = X.columns
        X = csr_matrix(X.values)

        val = self.transform(validate, "relevant")

        factors = 256
        num_epochs = 5
        self.batch_size = 32
        self.model = Elsa(data)
        self.model.fit_with_params(
            X=X,
            validation_data=val,
            factors=factors,
            num_epochs=num_epochs,
            batch_size=self.batch_size,
            shuffle=True,
            device=DEVICE,
        )

    def get(self, X, limit):
        X = self.transform(X, "relevant")
        y_pred = self.model.predict(X, self.batch_size)
        topk = torch.topk(y_pred, limit, sorted=True)
        return list(map(lambda indices: self.columns[indices].to_list(), topk.indices))

    def transform(self, X, column_name):
        ids = X.index
        X = X[[column_name]].explode(column_name)
        X = X[X[column_name].isin(self.columns)]
        X = pd.crosstab(X.index, X["relevant"])
        X = X.reindex(ids, columns=self.columns, fill_value=0)
        X = csr_matrix(X.values)
        return X


class ElsaPrevExp:
    """
    Results

    """

    def __init__(self, data: DataRepository, train, validate):
        self.name = "ElsaPrev"

        X = train[["prev"]].explode("prev")
        X = pd.crosstab(X.index, X["prev"])
        self.columns = X.columns
        X = csr_matrix(X.values)

        val = self.transform(validate, "prev")

        factors = 256
        num_epochs = 5
        self.batch_size = 32
        self.model = Elsa(data)
        self.model.fit_with_params(
            X=X,
            validation_data=val,
            factors=factors,
            num_epochs=num_epochs,
            batch_size=self.batch_size,
            shuffle=True,
            device=DEVICE,
        )

    def get(self, X, limit):
        X = self.transform(X, "prev")
        y_pred = self.model.predict(X, self.batch_size)
        topk = torch.topk(y_pred, limit, sorted=True)
        return list(map(lambda indices: self.columns[indices].to_list(), topk.indices))

    def transform(self, X, column_name):
        ids = X.index
        X = X[[column_name]].explode(column_name)
        X = X[X[column_name].isin(self.columns)]
        X = pd.crosstab(X.index, X["prev"])
        X = X.reindex(ids, columns=self.columns, fill_value=0)
        X = csr_matrix(X.values)
        return X
