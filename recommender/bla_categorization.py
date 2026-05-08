import numpy as np
import pandas as pd
from data import TrainData
from explainer.elsa import ElsaExplainer
from explainer.gcn import Categories
from grouper import RankCategorizer, SyntaxGrouper
from ranker.elsa_ranker import Elsa as ElsaRanker
from ranker.gcn import GCNRanker
from ranker.lightgcn import LightGCNRanker
from user import User


def main():
    train_data = TrainData(67564)
    train_data.fit()

    user = User(train_data.rand_soident_from_dev(), "", 0, None)
    user.fetch = True

    # ranker = GCNRanker(train_data)
    # ranker.set_train_params(epochs=5)
    # ranker.fit()
    ranker = ElsaRanker(train_data)
    ranker.fit()

    grouper = SyntaxGrouper(train_data)
    grouper.fit()

    pred = ranker.rank(user)
    groups = grouper.group(pred, 10)
    limit = sum([len(g) for g in groups])

    categorizer = RankCategorizer(train_data)
    cat_names, cat_values = categorizer.categorize(pred)
    cat_groups = [grouper.group(c) for c in cat_values]

    explainer = ElsaExplainer(ranker, train_data)
    explainer.fit()
    category_pred = [course for category in cat_values for course in category]
    expl = list(set(pred[:limit] + category_pred))
    print(
        "pred(",
        len(pred[:limit]),
        ") + cat(",
        len(category_pred),
        ") = ",
        len(expl),
        sep="",
    )
    expl = explainer.explain(user, list(expl))  # TODO: make dict

    for name, values, groups in zip(cat_names, cat_values, cat_groups):
        print(name)
        for g in groups:
            print(4 * " ", "- ", end="")
            for vi in g:
                c = values[vi]
                print(c, "(", expl[c], ")", end=", ", sep="")
            print()


if __name__ == "__main__":
    main()
