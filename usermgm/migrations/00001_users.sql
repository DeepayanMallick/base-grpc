-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(

    id                          VARCHAR(100)    PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name                  VARCHAR(100)    UNIQUE NOT NULL DEFAULT '',
    last_name                   VARCHAR(100)    UNIQUE NOT NULL DEFAULT '',
    username                    VARCHAR(100)    UNIQUE NOT NULL DEFAULT '',
    email                       VARCHAR(150)    UNIQUE NOT NULL DEFAULT '',
    password                    VARCHAR(255)    NOT NULL DEFAULT '',
    status                      SMALLINT        DEFAULT 0,
    created_at                  TIMESTAMP       DEFAULT current_timestamp,
    updated_at                  TIMESTAMP       DEFAULT current_timestamp,
    deleted_at                  TIMESTAMP       DEFAULT NULL,
    deleted_by                  VARCHAR(100)    NOT NULL DEFAULT ''
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;
