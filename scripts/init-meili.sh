#!/bin/bash

# Set the MeiliSearch API key
API_KEY=$MEILI_MASTER_KEY
BASE_URL="http://localhost:7700"

# Set filterable attributes
filterable='[
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
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$filterable" \
  "$BASE_URL/indexes/courses/settings/filterable-attributes")
echo "PUT filterable attributes: $response"

response=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "bert": {
        "source": "rest", 
        "url": "http://bert:8003/embedding", 
        "request": {
            "text": "{{text}}"
        },
        "response": {
            "embedding": "{{embedding}}"
        },
        "documentTemplate": "University course with title {{doc.title.en}}"
    }
  }' \
  "$BASE_URL/indexes/courses/settings/embedders")
echo "PATCH embedders: $response"

# Set searchable attributes
searchable='[
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
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$searchable" \
  "$BASE_URL/indexes/courses/settings/searchable-attributes")
echo "PUT searchable attributes: $response"

# Set dictionary
dict='[
  "C#",
  "c#",
  "C++",
  "c++"
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$dict" \
  "$BASE_URL/indexes/courses/settings/dictionary")
echo "PUT dictionary: $response"

filterable='[
    "teacher.id",
    "study_field.id",
    "academic_year",
    "study_year",
    "course_code",
    "study_type.id",
    "target_type"
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$filterable" \
  "$BASE_URL/indexes/survey/settings/filterable-attributes")
echo "PUT filterable attributes: $response"

sortable='[
    "academic_year",
    ""
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$sortable" \
  "$BASE_URL/indexes/survey/settings/sortable-attributes")
echo "PUT filterable attributes: $response"
