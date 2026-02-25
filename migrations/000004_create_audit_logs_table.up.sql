CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    event_type VARCHAR(64) NOT NULL,
    todo_id UUID,
    user_id UUID,
    payload TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at);

CREATE INDEX IF NOT EXISTS idx_audit_logs_event_type ON audit_logs (event_type);