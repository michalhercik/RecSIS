import random
import torch
from transformers import BertTokenizer, BertModel
from sklearn.metrics.pairwise import cosine_similarity
from flask import Flask, jsonify, request

RandomSeed = 52
random.seed(RandomSeed)
torch.manual_seed(RandomSeed)
if torch.cuda.is_available():
    torch.cuda.manual_seed_all(RandomSeed)

app = Flask(__name__)

@app.route("/embedding", methods=["POST"])
def embedding():
    r = request.get_json()
    text = r.get("text", "")
    embedding = bert_embed(text)
    response = jsonify({"embedding": embedding.flatten().tolist()})
    return response

tokenizer = BertTokenizer.from_pretrained('bert-base-uncased')
model = BertModel.from_pretrained('bert-base-uncased')

def bert_embed(text: str):
    encoding = tokenizer.batch_encode_plus(
        [text],                	
        padding=True,          	
        truncation=True,       	
        return_tensors='pt',  	
        add_special_tokens=True
    )
    
    token_ids = encoding['input_ids']  
    attentionMask = encoding['attention_mask']  

    with torch.no_grad():
        outputs = model(token_ids, attention_mask=attentionMask)
        word_embeddings = outputs.last_hidden_state  
        sentence_embedding = word_embeddings.mean(dim=1)

    return sentence_embedding

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8003, debug=True)
