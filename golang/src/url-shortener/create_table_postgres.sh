#!/bin/bash
if [ -z "$1" ]; then
    echo "$0 [Database name]"
    exit 1
fi

DBNAME=$1

set +x
psql --host="${HOST}" --port="${PORT}" --user="${USER}" -d "${DBNAME}" -c 'CREATE TABLE url (
    id VARCHAR(32) PRIMARY KEY,
    url VARCHAR(32)
);'

