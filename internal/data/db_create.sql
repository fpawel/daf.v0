PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS party
(
    party_id   INTEGER PRIMARY KEY NOT NULL,
    created_at TIMESTAMP           NOT NULL                     DEFAULT (datetime('now')) UNIQUE,
    type       INTEGER             NOT NULL                     DEFAULT 1 CHECK ( type > 0 ),
    pgs1       REAL                NOT NULL CHECK ( pgs1 >= 0 ) DEFAULT 0,
    pgs2       REAL                NOT NULL CHECK ( pgs2 >= 0 ) DEFAULT 4,
    pgs3       REAL                NOT NULL CHECK ( pgs3 >= 0 ) DEFAULT 7.5,
    pgs4       REAL                NOT NULL CHECK ( pgs4 >= 0 ) DEFAULT 12
);

CREATE TABLE IF NOT EXISTS product
(
    product_id INTEGER PRIMARY KEY NOT NULL,
    party_id   INTEGER             NOT NULL,
    created_at TIMESTAMP           NOT NULL DEFAULT (datetime('now')) UNIQUE,
    serial     INTEGER             NOT NULL CHECK (serial > 0 ),
    addr       SMALLINT            NOT NULL CHECK (addr > 0),
    checked    BOOLEAN             NOT NULL DEFAULT FALSE,

    UNIQUE (party_id, addr),
    UNIQUE (party_id, serial),
    FOREIGN KEY (party_id) REFERENCES party (party_id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS product_value
(
    product_id    INTEGER NOT NULL,
    tex           BOOLEAN NOT NULL CHECK ( tex IN (0, 1) ),
    gas           INTEGER NOT NULL CHECK ( tex IN (1, 2, 3, 4) ),
    concentration REAL    NOT NULL,
    current       REAL    NOT NULL,
    threshold1    BOOLEAN NOT NULL,
    threshold2    BOOLEAN NOT NULL,
    UNIQUE (product_id, gas, tex),
    FOREIGN KEY (product_id) REFERENCES product (product_id) ON DELETE CASCADE
);

CREATE VIEW IF NOT EXISTS last_party AS
SELECT *
FROM party
ORDER BY created_at DESC
LIMIT 1;