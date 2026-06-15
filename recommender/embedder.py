from sentence_transformers import SentenceTransformer, models

model = "sentence-transformers/all-MiniLM-L6-v2"

word_embedding_model = models.Transformer(model)
pooling_model = models.Pooling(
    word_embedding_model.get_word_embedding_dimension(), pooling_mode_mean_tokens=True
)

sbert = SentenceTransformer(modules=[word_embedding_model, pooling_model])


def sbert_embed(texts):
    embeddings = sbert.encode(texts, normalize_embeddings=True)
    return embeddings
