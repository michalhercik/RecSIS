import matplotlib.cm as cm
import matplotlib.pyplot as plt
import numpy as np
from data import TrainData
from ranker.elsa_ranker import Elsa
from sklearn.manifold import TSNE


def main():
    train_data = TrainData(42)
    train_data.fit()
    model = Elsa(train_data)
    model.fit()
    print(model.clusters.sort_values("cluster").head(50))

    embeddings = np.vstack(model.clusters["reduced_embedding"].values)
    labels = model.clusters["cluster"].values

    model.clusters["povinn"]

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


if __name__ == "__main__":
    main()
