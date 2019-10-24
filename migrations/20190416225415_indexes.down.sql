
ALTER TABLE increments DROP CONSTRAINT increments_client_guid_fk_idx;

ALTER TABLE clients DROP CONSTRAINT clients_pk_idx;
ALTER TABLE increments DROP CONSTRAINT increments_pk_idx;

DROP INDEX clients_username_uniq_idx;
