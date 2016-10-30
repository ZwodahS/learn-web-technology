
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

if [ -z "${PGDBNAME}" ]; then
    PGDBNAME=${PGDBNAME}
    exit 1
fi

HAS=$(psql -l -t | cut -d '|' -f 1 | grep "${PGDBNAME}" | wc -l)
CREATE=0
if [ "${HAS}" != "0" ]; then
    if [ -n "${DROP_DB}" ]; then
        echo "Dropping DB"
        psql -c "DROP DATABASE ${PGDBNAME}"
        CREATE=1
    fi
else
    CREATE=1
fi

if [ "${CREATE}" != "1" ]; then
    echo "Nothing created"
    exit 0
fi

psql -c "CREATE DATABASE ${PGDBNAME}"
psql -d "${PGDBNAME}" -c "CREATE EXTENSION pgcrypto"
psql -d "${PGDBNAME}" -c "CREATE EXTENSION citext"
psql -d "${PGDBNAME}" -c 'CREATE TABLE users (
    id UUID PRIMARY KEY,
    username CITEXT UNIQUE,
    email CITEXT UNIQUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);'
psql -d "${PGDBNAME}" -c 'CREATE TABLE bookmarks (
    id UUID PRIMARY KEY,
    user_id uuid,
    url VARCHAR(100),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);'
psql -d "${PGDBNAME}" -c 'CREATE TABLE bookmark_tags (
    bookmark_id UUID,
    tag VARCHAR(20),
    PRIMARY KEY(bookmark_id, tag)
);'
