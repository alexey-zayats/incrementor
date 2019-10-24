
ALTER TABLE clients ADD CONSTRAINT clients_pk_idx PRIMARY KEY (guid);
ALTER TABLE increments ADD CONSTRAINT increments_pk_idx PRIMARY KEY (guid);

CREATE UNIQUE INDEX clients_username_uniq_idx ON clients (username);

ALTER TABLE increments ADD CONSTRAINT increments_client_guid_fk_idx FOREIGN KEY (client_guid) REFERENCES clients(guid);
