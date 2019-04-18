
ALTER TABLE increments DROP CONSTRAINT increments_username_fk_idx;

ALTER TABLE clients DROP CONSTRAINT clients_pk_idx;
ALTER TABLE increments DROP CONSTRAINT increments_pk_idx;

DROP INDEX clients_username_uniq_idx;
DROP INDEX increments_username_uniq_idx;