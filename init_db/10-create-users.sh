#!/bin/bash

psql -U postgres -c "CREATE USER elt WITH PASSWORD '$ELT_PASS';"
psql -U postgres -c "CREATE USER webapp WITH PASSWORD '$WEBAPP_PASS';"

