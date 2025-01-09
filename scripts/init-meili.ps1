$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/documents" `
    -Method Delete `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" }

$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/documents?primaryKey=id" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/x-ndjson" `
    -InFile "$PSScriptRoot/../init_search/courses.json"

echo "$($response.StatusCode) $($response.StatusDescription)"
echo $response.Content