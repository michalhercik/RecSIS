#!/bin/bash

# Set the MeiliSearch API key
API_KEY="MASTER_KEY"
BASE_URL="http://localhost:7700"

# Delete all documents from the "courses" index
response=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE \
  -H "Authorization: Bearer $API_KEY" \
  "$BASE_URL/indexes/courses/documents")
echo "DELETE documents: $response"

# Add documents to the "courses" index
response=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/x-ndjson" \
  --data-binary @"$(dirname "$0")/../init_search/courses.json" \
  "$BASE_URL/indexes/courses/documents?primaryKey=id")
echo "POST documents: $response"

# Set filterable attributes
filterable='[
  "start_semester",
  "semester_count",
  "lecture_range",
  "seminar_range",
  "section",
  "credits",
  "department",
  "exam_type",
  "range_unit",
  "taught",
  "language",
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

# Set searchable attributes
searchable='[
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

response=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/x-ndjson" \
  --data-binary @"$(dirname "$0")/../init_search/comments.json" \
  "$BASE_URL/indexes/courses-comments/documents?primaryKey=id")
echo "POST documents: $response"

filterable='[
    "teacher_facet",
    "study_field",
    "academic_year",
    "study_year",
    "course_code",
    "study_type.code",
    "study_type.name_cs",
    "study_type.name_en",
    "target_type"
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$filterable" \
  "$BASE_URL/indexes/courses-comments/settings/filterable-attributes")
echo "PUT filterable attributes: $response"

sortable='[
    "academic_year",
    ""
]'
response=$(curl -s -o /dev/null -w "%{http_code}" -X PUT \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$sortable" \
  "$BASE_URL/indexes/courses-comments/settings/sortable-attributes")
echo "PUT filterable attributes: $response"

response=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/x-ndjson" \
  --data-binary @"$(dirname "$0")/../init_search/degree-plans-transformed.json" \
  "$BASE_URL/indexes/degree-plans/documents?primaryKey=id")
echo "POST documents: $response"