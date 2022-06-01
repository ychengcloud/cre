DROP TABLE IF EXISTS "numeric";
CREATE TABLE "numeric" (
    "smallint" smallint NOT NULL,
    "smallint1" smallint,
    "smallint2" smallint UNIQUE,
    "smallint3" smallint DEFAULT '10',
    "integer" integer,
    "bigint" bigint,
    "boolean" boolean,
    "numeric" numeric,
    "real" real,
    "double precision" double precision,
    "money" money
);

COMMENT ON COLUMN "public"."numeric"."smallint2" IS 'smallint comment';

DROP TABLE IF EXISTS "string";
CREATE TABLE "string" (
  "bit" bit NOT NULL,
  "bit varying" bit varying(255) NOT NULL,
  "character" character NOT NULL,
  "character varying" character varying NOT NULL,
  "character varying255" character varying(255) NOT NULL,
  "text" text NOT NULL,
  "tsquery" tsquery NOT NULL,
  "tsvector" tsvector NOT NULL,
  "uuid" uuid NOT NULL,
  "xml" xml NOT NULL,
  "json" json NOT NULL,
  "jsonb" jsonb NOT NULL
);

DROP TABLE IF EXISTS "time";
CREATE TABLE "time" (
  "date" date NOT NULL,
  "time" time NOT NULL,
  "timestamp" timestamp NOT NULL,
  "timestamptz" timestamptz NOT NULL,
  "interval" interval NOT NULL
);

DROP TABLE IF EXISTS "binary";
CREATE TABLE "binary" (
  "bytea" bytea NOT NULL
);

DROP TABLE IF EXISTS "spatial";
CREATE TABLE "spatial" (
  "cidr" cidr NOT NULL,
  "inet" inet NOT NULL,
  "macaddr" macaddr NOT NULL,
  "box" box NOT NULL,
  "circle" circle NOT NULL,
  "line" line NOT NULL,
  "lseg" lseg NOT NULL,
  "path" path NOT NULL,
  "point" point NOT NULL,
  "polygon" polygon NOT NULL
);

CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');

DROP TABLE IF EXISTS "enum";
CREATE TABLE "enum" (
  "enum" mood
);

DROP TABLE IF EXISTS "fk1";
CREATE TABLE "fk1" (
  "id" bigint NOT NULL,
  "fkid" bigint NOT NULL,
  PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "fk2";
CREATE TABLE "fk2" (
  "id" bigint NOT NULL,
  PRIMARY KEY ("id")
);

ALTER TABLE "fk1" ADD FOREIGN KEY ("fkid") REFERENCES "fk2" (id);