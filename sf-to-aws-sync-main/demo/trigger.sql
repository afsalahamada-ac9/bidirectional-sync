-- creates the function notify(), or else will replace the existing one
CREATE OR REPLACE FUNCTION notify()
-- specifies the return type, $$ is used to define the body of the function -> allows you to include single quotesi n your function without needing to escape them
RETURNS TRIGGER AS $$
BEGIN
    -- pg_notify function takes in the name of the channel to which the notification is sent
    IF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('events',TG_TABLE_NAME || ' ' || TG_OP || ' ' || OLD.id);
        RETURN OLD;
    ELSE
        PERFORM pg_notify('events',TG_TABLE_NAME || ' ' || TG_OP || ' ' || NEW.id);
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- CREATE TABLE table_names(
--     id SERIAL PRIMARY KEY,
--     data TEXT
-- );

CREATE TRIGGER trigger
AFTER INSERT OR UPDATE OR DELETE ON courses
FOR EACH ROW EXECUTE FUNCTION notify();