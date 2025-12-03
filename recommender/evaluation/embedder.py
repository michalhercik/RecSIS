import sys

sys.path.insert(0, "..")
from algo.embedder import Embedder, EmbedderAnnotation, EmbedderSyllabus
from data_repository import DataRepository


class EmbedderExp1:
    """
    Meilisearch Embedding
        - from course title
    Query
        - from course titles
        - with prefix: Give me recommendations for similar courses like:
    Filter
        - exclude blueprint courses
        - include only computer science section
    Results
        - done in ~6 min (1363 samples)
    """

    def __init__(self, data: DataRepository):
        self.name = "Embedder1"
        self.model = Embedder(data)
        self.model.fit()

    def get(self, X, incompatible, interchangable, limit):
        query = self.model.build_query(X)
        filter = self.model.build_filter(X)
        y_pred = self.model.fetch_similar(query, filter, limit)
        return y_pred


class EmbedderExp2:
    """
    Meilisearch Embedding
        - from course title
    Query
        - from course title
        - with prefix: Give me recommendations for similar courses like:
    Filter
        - exclude
            - blueprint courses
            - incompatible with
            - interchangable for
        - include only computer science section
    Results
        - done in ~6 min (1363 samples)
    """

    def __init__(self, data: DataRepository):
        self.name = "Embedder2"
        self.model = Embedder(data)
        self.model.fit()

    def get(self, X, incompatible, interchangable, limit):
        query = self.model.build_query(X)
        filter_out = []
        filter_out.extend(X)
        filter_out.extend(incompatible)
        filter_out.extend(interchangable)
        filter = self.model.build_filter(set(filter_out))
        y_pred = self.model.fetch_similar(query, filter, limit)
        return y_pred


class EmbedderExp3:
    """
    Meilisearch Embedding
        - from course title
    Query
        - from course titles
        - without prefix
    Filter
        - exclude
            - blueprint courses
            - incompatible with
            - interchangable for
        - include only computer science section
    Results
    - done in ~6 min (1363 samples)
    """

    def __init__(self, data: DataRepository):
        self.name = "Embedder3"
        self.model = Embedder(data)
        self.model.fit()

    def get(self, X, incompatible, interchangable, limit):
        query = self.model.build_query(X, prefix="")
        filter_out = []
        filter_out.extend(X)
        filter_out.extend(incompatible)
        filter_out.extend(interchangable)
        filter = self.model.build_filter(set(filter_out))
        y_pred = self.model.fetch_similar(query, filter, limit)
        return y_pred


class EmbedderAnnotationExp:
    """
    Meilisearch Embedding
        - from course
            - title
            - annotation
    Query
        - from course
            - titles
            - annotation
        - without prefix
    Filter
        - exclude
            - blueprint courses
            - incompatible with
            - interchangable for
        - include only computer science section
    Results
        - Computationaly demanding ~20 min (1363 samples)

    |          |precision    |precision    |recall    |recall    |avg_prec    |avg_prec    |
    |----------|-------------|-------------|----------|----------|------------|------------|
    | relevant |mean         |std          |mean      |std       |mean        |std         |
    |        k |             |             |          |          |            |            |
    |        10|0.0806       |0.1049       |0.0461    |0.0856    |0.1519      |0.2247      |
    |        20|0.0807       |0.0811       |0.0843    |0.1097    |0.1542      |0.1827      |
    |        50|0.0589       |0.0498       |0.148     |0.1413    |0.139       |0.1574      |
    |----------|-------------|-------------|----------|----------|------------|------------|
    |   next   |mean         |std          |mean      |std       |mean        |std         |
    |        k |             |             |          |          |            |            |
    |        10|0.0156       |0.043        |0.0323    |0.105     |0.0402      |0.1405      |
    |        20|0.0172       |0.0318       |0.0626    |0.1354    |0.0461      |0.1252      |
    |        50|0.0141       |0.0185       |0.121     |0.1804    |0.0474      |0.1137      |
    """

    def __init__(self, data: DataRepository):
        self.name = "EmbedderAnnotation"
        self.model = EmbedderAnnotation(data)
        self.model.fit()

    def get(self, X, incompatible, interchangable, limit):
        query = self.model.build_query(X, prefix="")
        filter_out = []
        filter_out.extend(X)
        filter_out.extend(incompatible)
        filter_out.extend(interchangable)
        filter = self.model.build_filter(set(filter_out))
        y_pred = self.model.fetch_similar(query, filter, limit)
        return y_pred


class EmbedderSyllabusExp:
    """
    Meilisearch Embedding
        - from course
            - title
            - syllabus
    Query
        - from course
            - title
            - syllabus
        - without prefix
    Filter
        - exclude
            - blueprint courses
            - incompatible with
            - interchangable for
        - include only computer science section
    Results
        - Computationaly demanding ~24 min (1363 samples)
    """

    def __init__(self, data: DataRepository):
        self.name = "EmbedderSyllabus"
        self.model = EmbedderSyllabus(data)
        self.model.fit()

    def get(self, X, incompatible, interchangable, limit):
        query = self.model.build_query(X, prefix="")
        filter_out = []
        filter_out.extend(X)
        filter_out.extend(incompatible)
        filter_out.extend(interchangable)
        filter = self.model.build_filter(set(filter_out))
        y_pred = self.model.fetch_similar(query, filter, limit)
        return y_pred
