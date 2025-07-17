$filterable = @(
    "start_semester",
    "semester_count",
    "lecture_range",
    "seminar_range",
    "section",
    "credits",
    "department",
    "exam",
    "range_unit",
    "taught_state",
    "taught_lang",
    "faculty",
    "capacity",
    "min_occupancy"
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/settings/filterable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($filterable | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"

$serchable = @(
    "code",
    "title",
    "guarantors",
    "teachers",
    "annotation",
    "sylabus",
    "aim",
    "terms_of_passing",
    "requirements_of_assesment",
    "literature"
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/settings/searchable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($searchable | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"

$dict = @("C#", "c#", "C++", "c++")
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/settings/dictionary" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($dict | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"

$filterable = @(
    "teacher.id",
    "study_field.id",
    "academic_year",
    "study_year",
    "course_code",
    "study_type.id",
    "target_type"
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/survey/settings/filterable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($filterable | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"

$sortable = @(
    "academic_year",
    ""
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/survey/settings/sortable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($sortable | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"

$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/degree-plans/documents?primaryKey=id" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/x-ndjson" `
    -InFile "$PSScriptRoot/../init_search/degree-plans-transformed.json"
echo "$($response.StatusCode) $($response.Content)"