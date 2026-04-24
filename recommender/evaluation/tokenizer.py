import re

import pandas as pd
from scipy.cluster.hierarchy import fcluster, leaves_list, linkage
from scipy.spatial.distance import squareform
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.metrics import pairwise_distances


def main():
    corpus = [
        "Anglický jazyk pro fyziky II",
        "Anglický jazyk pro informatiky II",
        "Anglický jazyk pro matematiky I",
        "Anglický jazyk pro středně pokročilé I",
        "Anglický jazyk pro mírně pokročilé I",
        # "Anglický jazyk pro mírně pokročilé III",
        # "Lineární algebra 1",
        # "Paralelní algoritmy",
        # "Lineární algebra 2",
        "Proseminář z lineární algebry",
        "Programování 1",
        "Programování 2",
        "Programování herních mechanik",
        "Programování v ChatGPT",
        "Programování v LabView pro fyziky",
        "Francouzský jazyk pro pokročilé I",
        "Německý jazyk pro pokročilé I",
        "Anglický jazyk pro pokročilé doktorandy II",
    ]
    # corpus = [
    #     "Robot 1",
    #     "Robot 2",
    #     "Transakce",
    #     "Algebra",
    #     "Forsing",
    #     "Algebra 1",
    #     "Algebra 2",
    #     "Algebra 2",
    #     "Algebra 1",
    #     "Algebra 2",
    #     "Účetnictví 1",
    #     "Ekonomie",
    #     "Programování 3",
    #     "Sage",
    #     "Psychologie",
    #     "Psychologie",
    #     "Lingvistika",
    #     "Visualizace",
    #     "Programování 1",
    #     "Programování 2",
    #     "Algoritmizace",
    #     "Programování I",
    #     "Middleware",
    #     "Middleware",
    #     "Složitost I",
    #     "Složitost",
    #     "Vyčíslitelnost",
    #     "Rekurze",
    # ]
    # vectorizer = CountVectorizer()
    # allow tokens that include trailing ++ or standalone ++

    # print(vectorizer.get_feature_names_out())
    # print(X.toarray())

    labels, Z, dist_matrix, vectorizer = assign_cluster_labels(
        corpus,
    )
    print(dist_matrix)
    df = pd.DataFrame({"name": corpus, "cluster": labels})
    print(df.sort_values("cluster"))
    print(vectorizer.get_feature_names_out())

    sim_matrix = 1 - dist_matrix
    df = pd.DataFrame(sim_matrix, index=corpus, columns=corpus)
    # # Basic heatmap
    # plot_heatmap(df, save_path="similarity_heatmap.png")
    # Clustered heatmap (optional)
    plot_clustered_heatmap(df, save_path="clustered_similarity_heatmap.png")


# -- Simple heatmap visualization -------------------------------------------
def plot_heatmap(df, figsize=(10, 8), cmap="YlGnBu", annot=True, save_path=None):
    plt.figure(figsize=figsize)
    sns.heatmap(
        df,
        annot=annot,
        fmt=".2f",
        cmap=cmap,
        square=True,
        cbar_kws={"label": "Jaccard similarity"},
    )
    plt.title("Jaccard similarity between course titles")
    plt.xticks(rotation=45, ha="right")
    plt.yticks(rotation=0)
    plt.tight_layout()
    if save_path:
        plt.savefig(save_path, dpi=200)
        print(f"Saved heatmap to: {save_path}")
    plt.show()


# -- Optional: reorder rows/cols by hierarchical clustering -----------------
def plot_clustered_heatmap(
    df,
    method="average",
    metric="euclidean",
    figsize=(10, 8),
    cmap="YlGnBu",
    annot=True,
    save_path=None,
):
    # linkage expects a condensed distance matrix or feature matrix. We'll cluster on the similarity rows.
    # Convert similarity to distance for clustering (higher distance = less similar)
    # Use 1 - similarity as a distance measure for clustering
    dist_for_clustering = 1 - df.values
    # Compute linkage on the rows
    Z = linkage(dist_for_clustering, method=method, metric=metric)
    order = leaves_list(Z)
    df_reordered = df.iloc[order, :].iloc[:, order]

    plt.figure(figsize=figsize)
    sns.heatmap(
        df_reordered,
        annot=annot,
        fmt=".2f",
        cmap=cmap,
        square=True,
        cbar_kws={"label": "Jaccard similarity"},
    )
    plt.title("Clustered Jaccard similarity (hierarchical reorder)")
    plt.xticks(rotation=45, ha="right")
    plt.yticks(rotation=0)
    plt.tight_layout()
    if save_path:
        plt.savefig(save_path, dpi=200)
        print(f"Saved clustered heatmap to: {save_path}")
    plt.show()


# remove arabic numbers and Roman numerals from course titles
roman_re = re.compile(
    r"\b(?:M{0,4}(?:CM|CD|D?C{0,3})(?:XC|XL|L?X{0,3})(?:IX|IV|V?I{0,3}))\b",
    flags=re.IGNORECASE,
)


def clean_title(s):
    s = re.sub(r"\d+", "", s)  # remove Arabic digits
    s = roman_re.sub("", s)  # remove Roman numerals
    # s = re.sub(r"[\(\)\.,:;/\-]+", " ", s)  # replace some punctuation with space
    # s = re.sub(r"\s+", " ", s).strip()  # collapse whitespace
    # if len(s) == 1 and s.lower() == "programování":
    #     return "_programování"
    # s = s.lower()
    stop_words = set(
        [
            "pro začátečníky",
            "pro mírně pokročilé",
            "pro středně pokročilé",
            "pro pokročilé",
            # "úvod",
            # "jazyk",
            # "programování",
            # "pro",
            # "z",
            # "v",
        ]
    )
    for w in stop_words:
        s = s.replace(w, "")
    s = re.sub(r"\s+", " ", s).strip()  # collapse whitespace
    return s


def equality_labels(corpus):
    norm = [clean_title(c) for c in corpus]
    vocab = {w: i for i, w in enumerate(set(norm))}
    clusters = [vocab[w] for w in norm]
    cluster_names = {cluster: name for name, cluster in vocab.items()}
    return clusters, cluster_names


def assign_cluster_labels(
    corpus,
    vectorizer=None,
    ngram_range=(2, 2),
    n_clusters=None,
    distance_threshold=None,
    linkage_method="average",
):
    """
    Assign a single integer cluster label to each item in `corpus` using hierarchical clustering.

    Parameters:
    - corpus: list of strings (course titles)
    - vectorizer: optional CountVectorizer instance (if None one will be created)
    - ngram_range: passed to CountVectorizer when vectorizer is None
    - n_clusters: if provided, produce this many clusters (uses criterion='maxclust')
    - distance_threshold: if provided (and n_clusters is None), cut the dendrogram at this distance
    - linkage_method: linkage method passed to scipy.cluster.hierarchy.linkage

    Returns:
    - labels: array-like of integer cluster labels (1..k)
    - Z: linkage matrix
    - dist_matrix: square pairwise distance matrix (Jaccard)
    - vectorizer: the fitted CountVectorizer
    """
    norm = [clean_title(c) for c in corpus]
    if vectorizer is None:
        token_pat = r"(?u)(?:\b\w+(?:\+\+)?\b|\+\+)"
        stop_words = [
            "mírně",
            "středně",
            "pokročilé",
            "začátečníky",
            "úvod",
            "jazyk",
            "programování",
            "pro",
            "z",
            "v",
        ]
        vectorizer = CountVectorizer(
            analyzer="word",
            ngram_range=(1, 2),
            token_pattern=token_pat,
            stop_words=stop_words,
        )
    X = vectorizer.fit_transform(norm)

    X_binary = (X > 0).astype(int)
    dist_matrix = pairwise_distances(X_binary.toarray(), metric="jaccard")

    # Convert to condensed distance matrix for linkage
    condensed = squareform(dist_matrix, checks=False)

    # Build linkage
    Z = linkage(condensed, method=linkage_method)

    # Decide how to form flat clusters
    if n_clusters is not None:
        labels = fcluster(Z, t=n_clusters, criterion="maxclust")
    elif distance_threshold is not None:
        labels = fcluster(Z, t=distance_threshold, criterion="distance")
    else:
        # default: cut at distance 0.5 (can be adjusted)
        labels = fcluster(Z, t=0.5, criterion="distance")

    return labels, Z, dist_matrix, vectorizer


def group_by_cluster(courses, clusters, k=None):
    print(courses.shape)
    print(courses[:k], flush=True)
    cluster_memory = {}
    result = []
    for cluster, course in zip(clusters, courses):
        print(cluster_memory, flush=True)
        if len(result) == k:
            break
        if cluster in cluster_memory:
            j = cluster_memory[cluster]
            result[j].append(course)
        else:
            result.append([course])
            cluster_memory[cluster] = len(result) - 1
    print(result)
    return result


if __name__ == "__main__":
    import matplotlib.pyplot as plt
    import seaborn as sns

    main()
