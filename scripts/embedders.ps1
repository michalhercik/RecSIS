$response = Invoke-RestMethod -Uri "http://localhost:7700/indexes/courses/settings/embedders" `
    -Method Get `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" 
echo $response