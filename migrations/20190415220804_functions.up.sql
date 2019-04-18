
CREATE OR REPLACE FUNCTION all_before_insert_update()
  RETURNS TRIGGER AS $$
BEGIN
  NEW.relname = TG_RELNAME;
  NEW.updated := now();
  IF TG_OP = 'UPDATE'
  THEN
    IF NEW.guid IS NULL
    THEN
      new.guid := old.guid;
    END IF;
  ELSEIF TG_OP = 'INSERT'
    THEN
      IF NEW.guid IS NULL
      THEN
        NEW.guid := uuid_generate_v4();
      END IF;
      IF NEW.created IS NULL
      THEN
        NEW.created := now();
      END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION all_before_insert_update() IS 'Updates the updated field, fills in the relname and monitors the created & guid';

CREATE OR REPLACE FUNCTION not_insert_update_delete()
  RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    RAISE EXCEPTION 'INSERT into table % not allowed', TG_RELNAME;
  ELSIF TG_OP = 'UPDATE' THEN
      RAISE EXCEPTION 'UPDATE of table % not allowed', TG_RELNAME;
  ELSEIF TG_OP = 'DELETE' THEN
      RAISE EXCEPTION 'DELETE of table % not allowed', TG_RELNAME;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION not_insert_update_delete() IS 'Forbids to insert, update and delete data directly in the base tables';
