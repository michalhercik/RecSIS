import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from adjustText import adjust_text
from data import TrainData
from ranker.elsa_ranker import Elsa
from sklearn.decomposition import PCA
from sklearn.manifold import TSNE
from sklearn.preprocessing import normalize

RND_STATE = 42


def main():
    train_data = TrainData(67564)
    train_data.fit()

    model = Elsa(train_data)
    model.set_train_params(
        factors=16, num_epochs=50, learning_rate=1e-2, batch_size=128
    )
    model.fit()

    emb_matrix = model.model.get_items_embeddings(as_numpy=True)
    labels = model.train_data.povinn.apply(
        lambda x: x["povinn"] + " | " + x["panazev"], axis=1
    ).values
    assert len(labels) == emb_matrix.shape[0]

    emb_matrix = normalize(emb_matrix, norm="l2")
    tsne = TSNE(n_components=2, random_state=RND_STATE, init="pca")
    emb2d = tsne.fit_transform(emb_matrix)

    sampled_label_names = [
        "NAIL029 | Machine Learning",
        # "NPFL129 | Introduction to Machine Learning with Python",
        # "NPRG043 | Recommended Programming Practices",
        "NOPT060 | Cooperative game theory seminar",
        # "NJAZ039 | Russian for Beginners I",
        # "NJAZ106 | Russian for Advanced Students I",
        "NPRG041 | Programming in C++",
        "NCGD003 | Gameplay Programming",
        "NDMI002 | Discrete Mathematics",
        "NDBI040 | Modern Database Systems",
        # "NJAZ170 | English for Advanced Students I",
        "NDBI048 | Data Science",
        "NDBI042 | Data Visualization Techniques",
        "NMMB532 | Standards and Cryptography",
    ]
    # sampled_label_names += list(
    #     np.random.RandomState(RND_STATE + 3).choice(labels, size=10, replace=False)
    # )
    plt.scatter(emb2d[:, 0], emb2d[:, 1], s=8, color="#cccccc", alpha=0.8)

    texts = []
    # color and annotate points belonging to each chosen label
    for idx, label_name in enumerate(sampled_label_names):
        mask = labels == label_name
        if not np.any(mask):
            continue
        color = "#305CDE"
        plt.scatter(
            emb2d[mask, 0],
            emb2d[mask, 1],
            s=20,
            color=color,
            label=str(label_name),
            alpha=0.9,
        )

        # annotate at the median position of that label's points
        centroid = np.median(emb2d[mask], axis=0)
        txt = plt.annotate(
            str(label_name),
            xy=(centroid[0], centroid[1]),
            xytext=(5, 5),
            textcoords="offset points",
            fontsize=11,
            color=color,
            bbox=dict(boxstyle="round,pad=0.2", fc="white", alpha=0.6, linewidth=0.5),
        )
        texts.append(txt)

    plt.tight_layout()
    plt.axis("off")
    plt.savefig("bla_clusters.png")
    plt.show()


if __name__ == "__main__":
    main()
