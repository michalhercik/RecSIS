import matplotlib.cm as cm
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from data import TrainData
from explainer.elsa import ElsaExplainer
from filterer import FinishedFilter
from ranker.elsa_ranker import Elsa
from sklearn.cluster import HDBSCAN, KMeans
from sklearn.decomposition import PCA
from sklearn.manifold import TSNE
from sklearn.preprocessing import normalize
from user import User


def main():
    train_data = TrainData(67564)
    train_data.fit()

    user = User(train_data.rand_soident_from_dev(), "", 0, None)
    user.fetch = True

    # print(train_data.klas["nazev"].value_counts())
    # print(train_data.klas["nazev"].drop_duplicates().shape)
    # print(train_data.trida["nazev"].value_counts().head(50))
    # print(train_data.trida["nazev"].drop_duplicates().shape)
    # exit()

    model = Elsa(train_data)
    model.set_train_params(
        factors=16, num_epochs=50, learning_rate=1e-2, batch_size=128
    )
    model.fit()

    # finished_filter = FinishedFilter(train_data)
    # explainer = ElsaExplain(model, train_data)

    # print(user.id)

    # pred = model.rank(user)
    # pred = finished_filter.filter(user, pred)
    # pred = pred[:10]
    # explain = explainer.explain(user, pred)
    # print(pd.DataFrame({"pred": pred, "explain": explain}))

    # results = (
    #     pd.DataFrame(model.rank(user), columns=["povinn"])
    #     .merge(train_data.povinn, on="povinn", how="left")
    #     .merge(
    #         train_data.klas.groupby("povinn")
    #         .agg({"nazev": list})
    #         .rename(columns={"nazev": "klas"}),
    #         on="povinn",
    #         how="left",
    #     )
    #     .merge(
    #         train_data.trida.groupby("povinn")
    #         .agg({"nazev": list})
    #         .rename(columns={"nazev": "trida"}),
    #         on="povinn",
    #         how="left",
    #     )[["povinn", "pgarant", "klas", "trida", "vucit1", "vucit2", "vucit3"]]
    # )

    clusters = cluster(
        model.model.get_items_embeddings(as_numpy=True),
        train_data,
        pca_components=50,  # Originally 256, cumulative explained variance is 0.37
        n_clusters=30,
        method="hdbscan",
        hdbscan_min_cluster_size=5,
        hdbscan_max_cluster_size=None,
        mode_tol=0,
    )

    summary = (
        clusters[
            [
                "cluster",
                "size",
                "trida_ratio",
                "trida",
                "klas_ratio",
                "klas",
                "pgarant_ratio",
                "pgarant",
                "categories_ratio",
                "categories",
                "ucit_ratio",
                "ucit",
            ]
        ]
        .drop_duplicates()
        .sort_values("size", ascending=False)
        .round(2)
    )
    print(
        summary[
            [
                "size",
                "trida",
                "categories",
                "ucit",
                "pgarant",
                "trida_ratio",
                "categories_ratio",
                "ucit_ratio",
                "pgarant_ratio",
            ]
        ]
    )
    print(
        summary[
            [
                "size",
                "trida_ratio",
                "klas_ratio",
                "pgarant_ratio",
                "categories_ratio",
                "ucit_ratio",
            ]
        ]
        .describe()
        .round(2)
    )
    # print(summary["trida"].value_counts())
    # print(summary["klas"].value_counts())
    # print(summary["pgarant"].value_counts())

    embeddings = np.vstack(clusters["reduced_embedding"].values)
    labels = (
        clusters.apply(
            lambda x: (
                str(x["cluster"])
                + " = ("
                + str(x["size"])
                + ") - ["
                + str(round((x["categories_ratio"]), 2))
                + "] - "
                + str(x["categories"])
            ),
            axis=1,
        )
        .astype(str)
        .values
    )

    visualize_embeddings_tsne(embeddings, labels)


def visualize_embeddings_tsne(
    embeddings, labels, save_path="embeddings_tsne.png", random_state=42
):
    """
    embeddings: (N, D) numpy array
    labels: length-N array-like of cluster ids (ints or strings)
    """
    if embeddings.shape[0] == 0:
        print("No embeddings to visualize.")
        return

    # Run t-SNE (use PCA init for stability)
    tsne = TSNE(n_components=2, random_state=random_state, init="pca")
    emb2d = tsne.fit_transform(embeddings)

    # Prepare colors for each cluster
    unique_labels = np.unique(labels)
    num_clusters = len(unique_labels)
    colormap = cm.get_cmap("tab10" if num_clusters <= 10 else "tab20", num_clusters)
    label_to_idx = {lab: i for i, lab in enumerate(unique_labels)}
    colors = [colormap(label_to_idx[l]) for l in labels]

    plt.figure(figsize=(8, 6))
    scatter = plt.scatter(emb2d[:, 0], emb2d[:, 1], c=colors, s=20, alpha=0.8)

    # Build a legend with one entry per cluster
    for lab in unique_labels:
        idx = label_to_idx[lab]
        plt.scatter([], [], color=colormap(idx), label=str(lab))
    plt.legend(title="cluster", bbox_to_anchor=(1.05, 1), loc="upper left")
    plt.title("t-SNE visualization of embeddings")
    plt.xlabel("tsne-1")
    plt.ylabel("tsne-2")
    plt.tight_layout()
    plt.savefig(save_path, dpi=150)
    print(f"Saved t-SNE plot to {save_path}")
    plt.show()


def cluster(
    emb_matrix,
    train_data,
    pca_components=50,
    n_clusters=None,
    method="kmeans",
    hdbscan_min_cluster_size=5,
    hdbscan_max_cluster_size=None,
    mode_tol=0.1,
):
    """
    Cluster embeddings using either KMeans or HDBSCAN.

    Parameters:
    - emb_matrix: numpy array of shape (n_items, dim)
    - train_data: TrainData object used to look up titles/classes
    - pca_components: components for PCA when original dim > 50
    - n_clusters: number of clusters for KMeans (ignored for HDBSCAN)
    - method: 'kmeans' (default) or 'hdbscan'
    - hdbscan_min_cluster_size: parameter passed to HDBSCAN
    - hdbscan_min_samples: parameter passed to HDBSCAN (can be None)
    """
    embeddings = pd.DataFrame(
        {
            "povinn": train_data.povinn["povinn"],
            "pnazev": train_data.povinn["pnazev"],
            "embedding": list(emb_matrix),
        }
    )
    emb_matrix = normalize(emb_matrix, norm="l2")

    orig_dim = emb_matrix.shape[1]
    if orig_dim > 50:
        pca = PCA(n_components=pca_components, random_state=42)
        reduced_emb = pca.fit_transform(emb_matrix)
        # print(pca.explained_variance_ratio_.cumsum())
    else:
        reduced_emb = emb_matrix.copy()

    # also keep reduced embeddings in the dataframe for later inspection
    embeddings["reduced_embedding"] = list(reduced_emb)

    # Build a matrix of embeddings (n_items x n_dims)
    # emb_matrix = np.vstack(embeddings["embedding"].values)
    n_items = emb_matrix.shape[0]

    method_lower = method.lower() if isinstance(method, str) else "kmeans"

    # Only compute a heuristic n_clusters when using kmeans and n_clusters is None
    if method_lower == "kmeans":
        # Heuristic for number of clusters: at least 1, otherwise sqrt(n_items)
        if n_clusters is None:
            if n_items <= 1:
                n_clusters = 1
            else:
                n_clusters = max(2, int(n_items**0.5))

    # Run clustering
    labels = None
    if method_lower == "hdbscan":
        clusterer = HDBSCAN(
            metric="cosine",
            min_cluster_size=hdbscan_min_cluster_size,
            max_cluster_size=hdbscan_max_cluster_size,
        )
        labels = clusterer.fit_predict(reduced_emb)

    else:
        # Default to k-means
        kmeans = KMeans(n_clusters=n_clusters, random_state=42, n_init=10)
        labels = kmeans.fit_predict(reduced_emb)

    # Attach cluster labels to the embeddings dataframe
    embeddings["cluster"] = labels

    cluster_size = embeddings.groupby("cluster")["povinn"].nunique()
    embeddings["size"] = embeddings["cluster"].map(cluster_size).fillna(0).astype(int)

    def test_title(title, column_name):
        merged = embeddings.merge(title, on="povinn", how="left")

        def get_modes(series):
            vc = series.dropna().value_counts()
            if vc.empty:
                return []
            maxc = vc.iloc[0]
            threshold = maxc * (1 - mode_tol)
            modes = vc[vc >= threshold].index.tolist()
            return modes

        mode_map = merged.groupby("cluster")["nazev"].agg(lambda s: get_modes(s))
        # mode_map = merged.groupby("cluster")["nazev"].agg(
        #     lambda x: x.mode().iat[0] if not x.mode().empty else np.nan
        # )

        # Attach the cluster_title (most frequent class name) to merged and embeddings
        merged["modes"] = merged["cluster"].map(mode_map)

        # Compute cluster size (number of items per cluster)
        cluster_size = merged.groupby("cluster")["povinn"].nunique()

        # Compute how many items in each cluster belong to the cluster_title
        # Use merged where we have the cluster_title value per row
        in_modes_mask = merged.apply(lambda row: row["nazev"] in row["modes"], axis=1)
        cluster_title_count = (
            # merged[merged["nazev"] == merged["trida"]]
            merged[in_modes_mask].groupby("cluster")["povinn"].nunique()
        )

        # Compute ratio: items in modal class / total items in cluster
        cluster_title_ratio = (cluster_title_count / cluster_size).fillna(0)

        # Map these cluster-level stats back to every item row in embeddings
        embeddings[column_name] = (
            embeddings["cluster"]
            .map(mode_map)
            .apply(lambda x: str(x) if len(x) != 1 else str(x[0]))
        )
        embeddings[f"{column_name}_ratio"] = (
            embeddings["cluster"].map(cluster_title_ratio).fillna(0)
        )

    test_title(train_data.trida, "trida")
    test_title(train_data.klas, "klas")
    pgarant = train_data.povinn[["povinn", "pgarant"]].copy()
    pgarant["nazev"] = pgarant["pgarant"]
    test_title(pgarant, "pgarant")
    test_title(train_data.ucit, "ucit")
    test_title(train_data.categories, "categories")

    return embeddings


if __name__ == "__main__":
    main()
