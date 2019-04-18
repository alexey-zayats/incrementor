
CREATE TABLE IF NOT EXISTS obj
(
  guid UUID NOT NULL DEFAULT uuid_generate_v4(),
  relname VARCHAR,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(guid)
);

COMMENT ON TABLE obj IS 'Base table for objects in database';
COMMENT ON COLUMN obj.guid IS 'Unique record identifier';
COMMENT ON COLUMN obj.relname IS 'Table name';
COMMENT ON COLUMN obj.created IS 'Date/time of record creation';
COMMENT ON COLUMN obj.updated IS 'Date/time of frecord last modification';

CREATE TABLE IF NOT EXISTS clients (
    username varchar not null,
    password varchar
) INHERITS (obj);

COMMENT ON TABLE clients IS 'Table holds records of clients that can query service';
COMMENT ON COLUMN clients.guid IS 'Unique record identifier';
COMMENT ON COLUMN clients.relname IS 'Table name';
COMMENT ON COLUMN clients.created IS 'Date/time of record creation';
COMMENT ON COLUMN clients.updated IS 'Date/time of frecord last modification';
COMMENT ON COLUMN clients.username IS 'Client username';
COMMENT ON COLUMN clients.password IS 'Client password';

CREATE TABLE IF NOT EXISTS increments (
    username varchar not null,
    number NUMERIC default 0,
    step NUMERIC default 0,
    maxvalue NUMERIC default 0
) INHERITS (obj);


COMMENT ON TABLE increments IS 'Table holds records of increments for clients';
COMMENT ON COLUMN increments.guid IS 'Unique record identifier';
COMMENT ON COLUMN increments.relname IS 'Table name';
COMMENT ON COLUMN increments.created IS 'Date/time of record creation';
COMMENT ON COLUMN increments.updated IS 'Date/time of frecord last modification';
COMMENT ON COLUMN increments.username IS 'Increment owner username';
COMMENT ON COLUMN increments.number IS 'Current number';
COMMENT ON COLUMN increments.step IS 'Step for increment';
COMMENT ON COLUMN increments.maxvalue IS 'Max increment value';