# #!/bin/bash

# MASTER_KEY="MASTER_KEY"
# COURSES_JSON_PATH="$(dirname "$0")/../init_search/courses.json"

# response=$(curl -s -X DELETE \
#     -H "Authorization: Bearer $MASTER_KEY" \
#     "http://localhost:7700/indexes/courses/documents")

# echo "$response"

# response=$(curl -s -X POST \
#     -H "Authorization: Bearer $MASTER_KEY" \
#     -H "Content-Type: application/x-ndjson" \
#     --data-binary "@$COURSES_JSON_PATH" \
#     "http://localhost:7700/indexes/courses/documents?primaryKey=id")

# echo "$response"

# dict='["C#", "c#", "C++", "c++"]'
# response=$(curl -s -X PUT \
#     -H "Authorization: Bearer $MASTER_KEY" \
#     -H "Content-Type: application/json" \
#     -d "$dict" \
#     "http://localhost:7700/indexes/courses/settings/dictionary")

# echo "$response"

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
  "lecture_range_winter",
  "seminar_range_winter",
  "lecture_range_summer",
  "seminar_range_summer",
  "credits",
  "faculty_guarantor",
  "exam_type",
  "range_unit",
  "taught",
  "taught_lang",
  "faculty",
  "capacity",
  "min_number"
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