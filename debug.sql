SET @id = 'test_stream_id';
INSERT INTO streams (stream_id, type, version)
SELECT @id, 'test', 0
WHERE NOT EXISTS (SELECT 1 FROM streams WHERE stream_id = @id AND version = 0);

-- Then proceed with your transaction
START TRANSACTION;
SET @stream_id = 'test_stream_id';
SET @event_id = 'test_event_id';
INSERT INTO events (event_id, stream_id, version, data, type)
SELECT @event_id, @stream_id, 1, '{"key": "value"}', 'test'
WHERE EXISTS (SELECT 1 FROM streams WHERE stream_id = @stream_id AND version = 0);
COMMIT;

-- Display the updated events table
SELECT * FROM streams;
SELECT * FROM events;

-- Display the updated streams table
SELECT * FROM streams WHERE id = @id AND version = 0;

DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS streams;