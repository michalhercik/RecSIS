#!/usr/bin/env bash
set -euo pipefail

# ----------------------------
# parameters
# ----------------------------
if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <migrations_dir> <container>"
    exit 1
fi

MIGRATIONS_DIR="$1"
CONTAINER="$2"

# ----------------------------
# configuration
# ----------------------------
DATABASE="${POSTGRES_DB:?POSTGRES_DB is not set}"
USER="${POSTGRES_USER:?POSTGRES_USER is not set}"

# ----------------------------
# validation
# ----------------------------
if [[ ! -d "$MIGRATIONS_DIR" ]]; then
    echo "Migration directory not found: $MIGRATIONS_DIR" >&2
    exit 1
fi

# find and sort SQL files
mapfile -t FILES < <(
    find "$MIGRATIONS_DIR" -maxdepth 1 -type f -name "*.sql" \
        | sort
)

if [[ ${#FILES[@]} -eq 0 ]]; then
    echo "No migration files found."
    exit 0
fi

echo "Running migrations in a single transaction..."

# ----------------------------
# build SQL stream
# ----------------------------
SQL=$(mktemp)

{
    echo "BEGIN;"
    echo "\set ON_ERROR_STOP on"

    for file in "${FILES[@]}"; do
        name=$(basename "$file")
        #echo "Including migration: $name"
        #echo "-- $name"
        cat "$file"
        echo
    done

    echo "COMMIT;"
} > "$SQL"

# ----------------------------
# execute once
# ----------------------------
docker exec -i "$CONTAINER" \
    psql \
    -U "$USER" \
    -d "$DATABASE" < "$SQL"

rm -f "$SQL"

echo "All migrations committed successfully."
