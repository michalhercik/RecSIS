$response = Invoke-RestMethod -Uri "http://localhost:7700/tasks/cancel?uids=$($args[0])" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json"
echo $response
