import torch
from sentence_transformers import SentenceTransformer, models

# consider Model2Vec https://github.com/MinishLab/model2vec?tab=readme-ov-file
word_embedding_model = models.Transformer("sentence-transformers/all-MiniLM-L12-v2")
pooling_model = models.Pooling(
    word_embedding_model.get_word_embedding_dimension(), pooling_mode_mean_tokens=True
)

sbert = SentenceTransformer(modules=[word_embedding_model, pooling_model])


def sbert_embed(texts):
    with torch.no_grad():
        embeddings = sbert.encode(texts, normalize_embeddings=True)
    return embeddings
