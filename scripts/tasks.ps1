$response = Invoke-RestMethod -Uri "http://localhost:7700/tasks" `
    -Method Get `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" 
echo $response.Results