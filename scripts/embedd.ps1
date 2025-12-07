$response = Invoke-RestMethod -Uri "http://localhost:8003/embedding" `
    -Method Post `
    -ContentType "application/json" `
    -Body '{
        "model": "sbert",
        "text": ["ahoj", "cau"]
    }'
echo $response
