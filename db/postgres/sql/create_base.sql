CREATE TABLE Link
(
    "base"  VARCHAR NOT NULL,
    "short" VARCHAR NOT NULL UNIQUE,
    CONSTRAINT "Link_pk" PRIMARY KEY ("base")
);
