import meilisearch
import json

client = meilisearch.Client('http://127.0.0.1:7700', 'MASTER_KEY')

index = client.index('courses')

with open('courses.json', 'r') as file:
    documents = file.read()

index.delete_all_documents()
index.add_documents_ndjson(documents.encode('utf-8'))