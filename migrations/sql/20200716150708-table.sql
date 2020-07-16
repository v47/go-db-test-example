
-- +migrate Up

-- +migrate StatementBegin
CREATE FUNCTION _doc_hash() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.doc_hash IS NULL THEN
        NEW.doc_hash := upper_md5(NEW.doc_name);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    doc_name TEXT NOT NULL,
    doc_hash VARCHAR(32) NOT NULL
);

CREATE TRIGGER user_set_doc_hash
BEFORE INSERT ON users
FOR EACH ROW
EXECUTE PROCEDURE _doc_hash();
-- +migrate StatementEnd


-- +migrate Down
