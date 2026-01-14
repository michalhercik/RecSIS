###############################################################
# Course index settings
###############################################################

$filterable = @(
    "code",
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
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($filterable | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"

$response = Invoke-RestMethod -Uri "http://localhost:7700/indexes/courses/settings/embedders" `
    -Method Patch `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body '{
        "bert": {
            "source": "rest", 
            "url": "http://bert:8003/embedding", 
            "request": {
                "text": "{{text}}"
            },
            "response": {
                "embedding": "{{embedding}}"
            },
            "documentTemplate": "University course with title {{doc.title.en}} {% if doc.annotation != nil %} {{doc.annotation[0]}} {% endif %}"
        }
    }'
Write-Output $response
# "documentTemplate": "University course with title {{doc.title.cs}} guaranted by {% for g in doc.guarantors %} {{ g.last_name }} {% endfor %} has following syllabus: {{doc.syllabus[0]}}"

$searchable = @(
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
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($searchable | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"

$dict = @("C#", "c#", "C++", "c++")
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/settings/dictionary" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($dict | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"

###############################################################
# Survey index settings
###############################################################

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
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($filterable | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"

$sortable = @(
    "academic_year",
    ""
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/survey/settings/sortable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($sortable | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"

###############################################################
# Degree plan index settings
###############################################################

$searchable = @(
    "code",
    "title"
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/degree-plans/settings/searchable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($searchable | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"

$filterable = @(
    "faculty",
    "section",
    "field.code",
    "teaching_lang",
    "validity",
    "study_type"
)
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/degree-plans/settings/filterable-attributes" `
    -Method Put `
    -Headers @{ "Authorization" = "Bearer $env:MEILI_MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($filterable | ConvertTo-Json)
Write-Output "$($response.StatusCode) $($response.Content)"
