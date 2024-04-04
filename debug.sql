use local;
select * from streams;
select * from events;

SET @uuid = '6ccd780c-baba-1026-9564-5b8c656024db';

INSERT INTO streams (id, type, version)
SELECT UUID_TO_BIN(@uuid), 'test', 0
WHERE NOT EXISTS (SELECT 1 FROM streams WHERE id = UUID_TO_BIN(@uuid) AND version = 0);

START TRANSACTION;
SET @uuid = '6ccd780c-baba-1026-9564-5b8c656024db';
INSERT INTO events (stream_id, version, data, type)
SELECT UUID_TO_BIN(@uuid), 1, '{"key": "value"}', 'test'
WHERE EXISTS (SELECT 1 FROM streams WHERE id = UUID_TO_BIN(@uuid) AND version = 0);
COMMIT;

SELECT * FROM streams;
SELECT * FROM events;


SELECT * FROM streams WHERE id = UUID_TO_BIN(@uuid) AND version = 0;

DROP TABLE streams;
DROP TABLE events;