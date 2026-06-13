import random

import torch
from flask import Flask, jsonify, request
from sentence_transformers import SentenceTransformer, models
from transformers import (
    BertModel,
    BertTokenizer,
)

app = Flask(__name__)

@app.route("/embedding", methods=["POST"])
def embedding():
    r = request.get_json()
    text = r.get("text", "")
    input = [text] if isinstance(text, str) else text
    embedding = sbert_embed(input)
    if isinstance(text, str):
        embedding = embedding.flatten()
    response = jsonify({"embedding": embedding.tolist()})
    return response

word_embedding_model = models.Transformer(
    "sentence-transformers/all-MiniLM-L6-v2"
)
pooling_model = models.Pooling(
    word_embedding_model.get_word_embedding_dimension(), pooling_mode_mean_tokens=True
)
sbert = SentenceTransformer(modules=[word_embedding_model, pooling_model])
def sbert_embed(texts):
    with torch.no_grad():
        embeddings = sbert.encode(
            texts, convert_to_tensor=True, normalize_embeddings=True
        )
    return embeddings


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8003, debug=True)
