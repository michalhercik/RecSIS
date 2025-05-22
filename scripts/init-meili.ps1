$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/documents" `
    -Method Delete `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" }
echo "$($response.StatusCode) $($response.Content)"

$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses-comments/documents" `
    -Method Delete `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" }
echo "$($response.StatusCode) $($response.Content)"

$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/documents?primaryKey=id" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/x-ndjson" `
    -InFile "$PSScriptRoot/../init_search/courses.json"
echo "$($response.StatusCode) $($response.Content)"

$filterable = @(
    "start_semester",
    "semester_count",
    "lecture_range",
    "seminar_range",
    # "lecture_range_winter",
    # "seminar_range_winter",
    # "lecture_range_summer",
    # "seminar_range_summer",
    "credits",
    "department",
    "exam_type",
    "range_unit",
    "taught",
    "language",
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
    "cs.name",
    "en.name",
    "guarantors",
    "teachers",
    "cs.A",
    "en.A",
    "cs.S",
    "en.S",
    "cs.C",
    "en.C",
    "cs.E",
    "en.E",
    "cs.P",
    "en.P",
    "cs.L",
    "en.L"
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

$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses-comments/documents?primaryKey=id" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/x-ndjson" `
    -InFile "$PSScriptRoot/../init_search/comments.json"
echo "$($response.StatusCode) $($response.Content)"

$filterable = @(
    "teacher_facet",
    "study_field",
    "academic_year",
    "study_year",
    "course_code",
    "study_type.code",
    "study_type.name_cs",
    "study_type.name_en",
    "target_type"
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses-comments/settings/filterable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($filterable | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"

$sortable = @(
    "academic_year",
    ""
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses-comments/settings/sortable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($sortable | ConvertTo-Json)
echo "$($response.StatusCode) $($response.Content)"