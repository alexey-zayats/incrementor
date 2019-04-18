\c template1

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE ROLE incrementor WITH LOGIN ENCRYPTED PASSWORD 'incrementor';
CREATE DATABASE incrementor WITH OWNER incrementor;