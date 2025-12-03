$response = Invoke-RestMethod -Uri "http://localhost:7700/indexes/courses/settings/embedders" `
    -Method Patch `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body '{
        "bert": {
            "source": "rest",
            "url": "http://bert:8003/embedding",
            "request": {
                "text": "{{text}}"
            },
            "response": {
                "embedding": "{{embedding}}"
            },
            "documentTemplate": "{{doc.embed_input}} "
        }
    }'
echo $response
