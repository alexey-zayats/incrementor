ALTER TABLE clients ADD CONSTRAINT clients_pk_idx PRIMARY KEY (guid);
ALTER TABLE increments ADD CONSTRAINT increments_pk_idx PRIMARY KEY (guid);

CREATE UNIQUE INDEX clients_username_uniq_idx ON clients (username);
CREATE UNIQUE INDEX increments_username_uniq_idx ON increments (username);

ALTER TABLE increments ADD CONSTRAINT increments_username_fk_idx FOREIGN KEY (username) REFERENCES clients(username);