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
        - BERT done in ~6 min (1363 samples)
        - SBERT done in ~1,5 min (1363 samples)
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
    |k         |             |             |          |          |            |            |
    |10        |0.0806       |0.1049       |0.0461    |0.0856    |0.1519      |0.2247      |
    |20        |0.0807       |0.0811       |0.0843    |0.1097    |0.1542      |0.1827      |
    |50        |0.0589       |0.0498       |0.148     |0.1413    |0.139       |0.1574      |
    |----------|-------------|-------------|----------|----------|------------|------------|
    |   next   |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.0156       |0.043        |0.0323    |0.105     |0.0402      |0.1405      |
    |20        |0.0172       |0.0318       |0.0626    |0.1354    |0.0461      |0.1252      |
    |50        |0.0141       |0.0185       |0.121     |0.1804    |0.0474      |0.1137      |
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
        - BERT - Computationaly demanding ~22 min (1363 samples)
        - SBERT - ~4 min (1363 samples)
        - BERT and SBERT similar results but SBERT is singificantly faster
        - We either don't have good embeddings or the students are not neccessarly interested in similar courses

    |   ALL    |precision    |precision    |recall    |recall    |avg_prec    |avg_prec    |
    |----------|-------------|-------------|----------|----------|------------|------------|
    | relevant |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.0894       |0.1022       |0.0477    |0.0870    |0.3321      |0.4011      |
    |20        |0.0718       |0.0706       |0.0785    |0.1111    |0.2857      |0.3308      |
    |50        |0.0580       |0.0475       |0.1606    |0.1566    |0.2025      |0.2138      |
    |----------|-------------|-------------|----------|----------|------------|------------|
    |   next   |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.0398       |0.0612       |0.0637    |0.1207    |0.2309      |0.3889      |
    |20        |0.0269       |0.0390       |0.0892    |0.1522    |0.2171      |0.3617      |
    |50        |0.0194       |0.0231       |0.1622    |0.2078    |0.1674      |0.2771      |


    |    BC    |precision    |precision    |recall    |recall    |avg_prec    |avg_prec    |
    |----------|-------------|-------------|----------|----------|------------|------------|
    | relevant |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.1005       |0.1054       |0.0427    |0.0671    |0.3882      |0.4171      |
    |20        |0.0804       |0.0715       |0.0709    |0.0897    |0.3312      |0.3438      |
    |50        |0.0656       |0.0487       |0.1473    |0.1247    |0.2306      |0.2217      |
    |----------|-------------|-------------|----------|----------|------------|------------|
    |   next   |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.0414       |0.0613       |0.0607    |0.1076    |0.2663      |0.4156      |
    |20        |0.0267       |0.0381       |0.0802    |0.1335    |0.2485      |0.3875      |
    |50        |0.0191       |0.0230       |0.1450    |0.1844    |0.1896      |0.2986      |


    |   MGR    |precision    |precision    |recall    |recall    |avg_prec    |avg_prec    |
    |----------|-------------|-------------|----------|----------|------------|------------|
    | relevant |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.0502       |0.0826       |0.0679    |0.1388    |0.1197      |0.2284      |
    |20        |0.0418       |0.0592       |0.1104    |0.1675    |0.1144      |0.1940      |
    |50        |0.0329       |0.0352       |0.2227    |0.2345    |0.0983      |0.1343      |
    |----------|-------------|-------------|----------|----------|------------|------------|
    |   next   |mean         |std          |mean      |std       |mean        |std         |
    |k         |             |             |          |          |            |            |
    |10        |0.0345       |0.0617       |0.0775    |0.1645    |0.0948      |0.2107      |
    |20        |0.0285       |0.0424       |0.1273    |0.2080    |0.0974      |0.1945      |
    |50        |0.0211       |0.0239       |0.2353    |0.2703    |0.0838      |0.1408      |
    """

    def __init__(self, data: DataRepository):
        self.name = "EmbedderSyllabus"
        self.model = EmbedderSyllabus(data)
        self.model.fit()

    def get(self, X, incompatible, interchangable, limit, prefix=""):
        query = self.model.build_query(X, prefix=prefix)
        filter_out = []
        filter_out.extend(X)
        filter_out.extend(incompatible)
        filter_out.extend(interchangable)
        filter = self.model.build_filter(set(filter_out))
        y_pred = self.model.fetch_similar(query, filter, limit)
        return y_pred
