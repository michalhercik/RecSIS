import time

from data import TrainData
from explainer.gcn import ExplainerType, GCNExplainer
from ranker.gcn import GCNRanker
from user import User


def main():
    train_data = TrainData(67564)
    train_data.fit()

    user = User(train_data.rand_soident_from_dev(), "", 0, None)
    user.fetch = True

    ranker = GCNRanker(train_data)
    ranker.set_train_params(epochs=50)
    ranker.fit()
    # pg_explainer = GCNExplainer(ranker, train_data)
    # pg_explainer.set_train_params(epochs=20, type=ExplainerType.PG)
    # pg_explainer.fit()
    # gcn_explainer = GCNExplainer(ranker, train_data)
    # gcn_explainer.set_train_params(epochs=20, type=ExplainerType.GNN)
    # gcn_explainer.fit()
    captum_explainer = GCNExplainer(ranker, train_data)
    captum_explainer.set_train_params(type=ExplainerType.CAPTUM)
    captum_explainer.fit()

    pred = ranker.rank(user)

    pred = pred[:10]
    print(pred)

    start = time.perf_counter()
    explanation = captum_explainer.explain(user, pred)
    elapsed = time.perf_counter() - start

    print("Explanation:", explanation)
    print(f"Explainer execution time: {elapsed:.6f} seconds")


if __name__ == "__main__":
    main()
