#!/bin/bash
rm -f url.db
sqlite3 url.db 'CREATE TABLE `url` (
    id VARCHAR(32) PRIMARY KEY,
    url VARCHAR(32)
);'

