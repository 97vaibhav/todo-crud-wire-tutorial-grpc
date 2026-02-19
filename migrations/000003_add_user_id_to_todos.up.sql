-- Existing todos (from dev testing) have no owner — clear them first.
-- In a real migration with live data you would assign a default user_id instead.
TRUNCATE TABLE todos;

ALTER TABLE todos
    ADD COLUMN user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;
