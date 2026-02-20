import random

import torch
from flask import Flask, jsonify, request
from sentence_transformers import SentenceTransformer, models
from sklearn.metrics.pairwise import cosine_similarity
from transformers import (
    BertModel,
    BertTokenizer,
    DPRContextEncoder,
    DPRContextEncoderTokenizer,
    DPRQuestionEncoder,
    DPRQuestionEncoderTokenizer,
)

RandomSeed = 52
random.seed(RandomSeed)
torch.manual_seed(RandomSeed)
if torch.cuda.is_available():
    torch.cuda.manual_seed_all(RandomSeed)

app = Flask(__name__)

SBERT = "sbert"
BERT = "bert"
DPR = "DPR"
DPR_QUERY = "query"
DPR_CONTEXT = "context"


@app.route("/embedding", methods=["POST"])
def embedding():
    r = request.get_json()
    model = r.get("model", SBERT)
    text = r.get("text", "")
    input = [text] if isinstance(text, str) else text
    if model == SBERT:
        embedding = sbert_embed(input)
    elif model == BERT:
        embedding = bert_embed(input)
    # elif model == DPR:
    #     mode = r.get("mode", "")
    #     if mode == DPR_QUERY:
    #         embedding = encode_query(input)
    #     elif mode == DPR_CONTEXT:
    #         embedding = encode_contexts(input)
    #     else:
    #         response = jsonify({"error": f"Unkown or missing mode: '{model}'"})
    #         response.status_code = 400
    #         return response
    else:
        response = jsonify({"error": f"Unkown model: '{model}'"})
        response.status_code = 400
        return response
    # print(embedding.shape, flush=True)
    if isinstance(text, str):
        embedding = embedding.flatten()
    response = jsonify({"embedding": embedding.tolist()})
    return response


tokenizer = BertTokenizer.from_pretrained("bert-base-uncased")
model = BertModel.from_pretrained("bert-base-uncased")


def bert_embed(texts):
    encoding = tokenizer.batch_encode_plus(
        texts,
        padding=True,
        truncation=True,
        return_tensors="pt",
        add_special_tokens=True,
    )

    token_ids = encoding["input_ids"]
    attentionMask = encoding["attention_mask"]

    with torch.no_grad():
        outputs = model(token_ids, attention_mask=attentionMask)
        word_embeddings = outputs.last_hidden_state
        sentence_embedding = word_embeddings.mean(dim=1)

    return sentence_embedding


# Explicit SBERT architecture (tokenizer + transformer + pooling)

word_embedding_model = models.Transformer(
    "sentence-transformers/all-MiniLM-L6-v2"
)  # super fast
# word_embedding_model = models.Transformer(
#     "sentence-transformers/all-mpnet-base-v2"
# )  # crashes - ~8GB not enough memory (exit 247)? maybe just because running batches
# word_embedding_model = models.Transformer("sentence-transformers/all-distilroberta-v1")
# word_embedding_model = models.Transformer(
#     "sentence-transformers/all-MiniLM-L12-v2"
# )  # super fast
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


# ----------------------------
# Load DPR models & tokenizers
# ----------------------------
# query_encoder_name = "facebook/dpr-question_encoder-single-nq-base"
# context_encoder_name = "facebook/dpr-ctx_encoder-single-nq-base"

# query_tokenizer = DPRQuestionEncoderTokenizer.from_pretrained(query_encoder_name)
# query_encoder = DPRQuestionEncoder.from_pretrained(query_encoder_name)

# context_tokenizer = DPRContextEncoderTokenizer.from_pretrained(context_encoder_name)
# context_encoder = DPRContextEncoder.from_pretrained(context_encoder_name)


# def encode_query(query):
#     inputs = query_tokenizer(query, return_tensors="pt", padding=True, truncation=True)
#     with torch.no_grad():
#         query_embedding = query_encoder(**inputs).pooler_output
#     return query_embedding


# def encode_contexts(passages):
#     inputs = context_tokenizer(
#         passages, return_tensors="pt", padding=True, truncation=True
#     )
#     with torch.no_grad():
#         context_embeddings = context_encoder(**inputs).pooler_output
#     return context_embeddings


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8003, debug=True)
