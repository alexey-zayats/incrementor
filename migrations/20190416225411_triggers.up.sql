
CREATE TRIGGER obj_before_insert_update BEFORE INSERT OR UPDATE OR DELETE ON obj FOR EACH ROW EXECUTE PROCEDURE not_insert_update_delete();
CREATE TRIGGER clients_before_insert_update BEFORE INSERT OR UPDATE ON clients FOR EACH ROW EXECUTE PROCEDURE all_before_insert_update();
CREATE TRIGGER increments_before_insert_update BEFORE INSERT OR UPDATE ON increments FOR EACH ROW EXECUTE PROCEDURE all_before_insert_update();
