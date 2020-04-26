BEGIN;
CREATE TABLE IF NOT EXISTS photos(
     id  SERIAL PRIMARY KEY,
     guid    int not null,
     name          varchar(255)      NOT NULL,
     path          varchar(400)       NOT NULL,
     vector        float8[] DEFAULT NULL
);
CREATE UNIQUE INDEX ON photos(path);
COMMIT;