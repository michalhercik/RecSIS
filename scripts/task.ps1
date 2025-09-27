$response = Invoke-RestMethod -Uri "http://localhost:7700/tasks/$($args[0])" `
    -Method Get `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" 
echo $response