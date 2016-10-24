
if [ -z "${PGUSER}" ]; then
    echo 'Set $PGUSER'
    exit 1
fi

if [ -z "${PGPASSWORD}" ]; then
    echo 'Set $PGPASSWORD'
    exit 1
fi

if [ -z "${PGHOST}" ]; then
    echo 'Set $PGHOST'
    exit 1
fi

if [ -z "${PGPORT}" ]; then
    echo 'Set $PGPORT'
    exit 1
fi

set +x
psql -c "DROP DATABASE IF EXISTS test_uuid"
psql -c "CREATE DATABASE test_uuid"
psql -d "test_uuid" -c "CREATE EXTENSION pgcrypto"
psql -d "test_uuid" -c "CREATE TABLE books (id uuid primary key, title varchar(50))"
