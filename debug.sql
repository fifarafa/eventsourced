use local;
select * from streams;
select * from events;

INSERT INTO streams (id, type, version)
SELECT 100, 'test', 0
WHERE NOT EXISTS (SELECT 1 FROM streams WHERE id = 100 AND version = 0);

INSERT INTO events (id, stream_id, version, data, type)
SELECT 200, 100, 1, '{}', 'test'
WHERE NOT EXISTS (SELECT 1 FROM streams WHERE id = 100 AND version = 0);

INSERT INTO events (id, stream_id, version, data, type)
SELECT UUID_TO_BIN(UUID()), 100, 1, '{"key": "value"}', 'test'
WHERE NOT EXISTS (SELECT 1 FROM streams WHERE id = 100 AND version = 0);