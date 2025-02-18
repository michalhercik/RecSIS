#!/bin/bash

MASTER_KEY="MASTER_KEY"
COURSES_JSON_PATH="$(dirname "$0")/../init_search/courses.json"

response=$(curl -s -X DELETE \
    -H "Authorization: Bearer $MASTER_KEY" \
    "http://localhost:7700/indexes/courses/documents")

echo "$response"

response=$(curl -s -X POST \
    -H "Authorization: Bearer $MASTER_KEY" \
    -H "Content-Type: application/x-ndjson" \
    --data-binary "@$COURSES_JSON_PATH" \
    "http://localhost:7700/indexes/courses/documents?primaryKey=id")

echo "$response"

dict='["C#", "c#", "C++", "c++"]'
response=$(curl -s -X PUT \
    -H "Authorization: Bearer $MASTER_KEY" \
    -H "Content-Type: application/json" \
    -d "$dict" \
    "http://localhost:7700/indexes/courses/settings/dictionary")

echo "$response"