-- +up
CREATE TABLE users (
    id           serial PRIMARY KEY,
    username     text,
    email        text,
		passwordhash text,

    elo          int,

		pwreset      boolean,
		confirmed    boolean
);

CREATE TABLE matches (
    id serial PRIMARY KEY
);

CREATE TABLE tiles (
    id         serial PRIMARY KEY,
    name       text,
    resistance int,
    defense    int,
    dodge      int
);

CREATE TABLE classes (
    id serial PRIMARY KEY,

    name text




);

-- +down
DROP TABLE users   CASCADE;
DROP TABLE matches CASCADE;
DROP TABLE tiles   CASCADE;
DROP TABLE classes CASCADE;
