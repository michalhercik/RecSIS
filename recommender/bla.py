import matplotlib.cm as cm
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from data import TrainData
from ranker.elsa_ranker import Elsa
from sklearn.cluster import KMeans
from sklearn.decomposition import PCA
from sklearn.manifold import TSNE
from sklearn.preprocessing import normalize


def main():
    train_data = TrainData(42)
    train_data.fit()
    model = Elsa(train_data)
    model.fit()

    clusters = cluster(
        model.model.get_items_embeddings(as_numpy=True),
        train_data,
        pca_components=50,  # Originally 256, cumulative explained variance is 0.37
        n_clusters=30,
    )

    summary = (
        clusters[["cluster", "size", "trida_ratio", "trida", "klas_ratio", "klas"]]
        .drop_duplicates()
        .sort_values("size", ascending=False)
        .round(2)
    )
    summary["ratio_diff"] = summary["trida_ratio"] - summary["klas_ratio"]
    print(summary)
    print(
        summary[["size", "trida_ratio", "klas_ratio", "ratio_diff"]].describe().round(2)
    )
    print(summary["trida"].value_counts())
    print(summary["klas"].value_counts())

    embeddings = np.vstack(clusters["reduced_embedding"].values)
    labels = (
        clusters.apply(
            lambda x: (
                str(x["cluster"])
                + " = ("
                + str(x["size"])
                + ") - ["
                + str(round((x["klas_ratio"] * 100), 2))
                + "] - "
                + str(x["klas"])
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


def cluster(emb_matrix, train_data, pca_components=50, n_clusters=None):
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
        print(pca.explained_variance_ratio_.cumsum())
    else:
        reduced_emb = emb_matrix.copy()

    # also keep reduced embeddings in the dataframe for later inspection
    embeddings["reduced_embedding"] = list(reduced_emb)

    # Build a matrix of embeddings (n_items x n_dims)
    # emb_matrix = np.vstack(embeddings["embedding"].values)
    n_items = emb_matrix.shape[0]

    # Heuristic for number of clusters: at least 1, otherwise sqrt(n_items)
    if n_clusters is None:
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

    cluster_size = embeddings.groupby("cluster")["povinn"].nunique()
    embeddings["size"] = embeddings["cluster"].map(cluster_size).fillna(0).astype(int)

    def test_title(title, column_name):
        merged = embeddings.merge(title, on="povinn", how="left")
        mode_map = merged.groupby("cluster")["nazev"].agg(
            lambda x: x.mode().iat[0] if not x.mode().empty else np.nan
        )
        # Attach the cluster_title (most frequent class name) to merged and embeddings
        merged["trida"] = merged["cluster"].map(mode_map)

        # Compute cluster size (number of items per cluster)
        cluster_size = merged.groupby("cluster")["povinn"].nunique()

        # Compute how many items in each cluster belong to the cluster_title
        # Use merged where we have the cluster_title value per row
        cluster_title_count = (
            merged[merged["nazev"] == merged["trida"]]
            .groupby("cluster")["povinn"]
            .count()
        )

        # Compute ratio: items in modal class / total items in cluster
        cluster_title_ratio = (cluster_title_count / cluster_size).fillna(0)

        # Map these cluster-level stats back to every item row in embeddings
        embeddings[column_name] = embeddings["cluster"].map(mode_map)
        embeddings[f"{column_name}_ratio"] = (
            embeddings["cluster"].map(cluster_title_ratio).fillna(0)
        )

    test_title(
        train_data.trida[train_data.trida["nazev"] != "Informatika Bc."], "trida"
    )
    test_title(
        train_data.klas[train_data.klas["nazev"] != "Předměty obecného základu"], "klas"
    )

    return embeddings


if __name__ == "__main__":
    main()
