$response = Invoke-WebRequest -Uri ("http://localhost:7700/tasks/" + $args[0]) `
    -Method GET `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" }

echo $response.Content