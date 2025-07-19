#!/bin/bash

psql -U postgres -c "CREATE USER elt WITH PASSWORD '$RECSIS_ELT_DB_PASS';"
psql -U postgres -c "CREATE USER webapp WITH PASSWORD '$RECSIS_WEBAPP_DB_PASS';"
psql -U postgres -c "CREATE USER recommender WITH PASSWORD '$RECSIS_RECOMMENDER_DB_PASS';"

