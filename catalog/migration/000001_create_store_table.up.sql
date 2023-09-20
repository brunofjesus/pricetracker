BEGIN;
CREATE TABLE IF NOT EXISTS STORE
(
    STORE_ID BIGSERIAL PRIMARY KEY,
    SLUG VARCHAR(255) NOT NULL,
    NAME VARCHAR(255) NOT NULL,
    WEBSITE VARCHAR(255) NOT NULL,
    ACTIVE BOOL NOT NULL DEFAULT TRUE,

    CONSTRAINT UK_STORE_SLUG UNIQUE (SLUG),
    CONSTRAINT UK_STORE_WEBSITE UNIQUE (WEBSITE)
);
COMMIT;