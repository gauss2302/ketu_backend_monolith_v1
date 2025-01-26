#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE ke2;
    GRANT ALL PRIVILEGES ON DATABASE ke2 TO postgres;
EOSQL 