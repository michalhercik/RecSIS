$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/documents" `
    -Method Delete `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" }
echo "$($response.StatusCode) $($response.StatusDescription)"
echo "  $($response.Content)"

$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/documents?primaryKey=id" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/x-ndjson" `
    -InFile "$PSScriptRoot/../init_search/courses.json"
echo "$($response.StatusCode) $($response.StatusDescription)"
echo "  $($response.Content)"

$dict = @("C#", "c#", "C++", "c++")
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/settings/dictionary" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($dict | ConvertTo-Json)
echo "$($response.StatusCode) $($response.StatusDescription)"
echo "  $($response.Content)"