
-- +migrate Up

-- +migrate StatementBegin
CREATE FUNCTION upper_md5(val text) RETURNS text AS $$
BEGIN
RETURN upper(md5(val));
END; $$
LANGUAGE PLPGSQL;
-- +migrate StatementEnd

-- +migrate Down

DROP FUNCTION upper_md5;
