-- Enable uuid-ossp extension so we can use uuid_generate_v4() as a default.
-- uuid-ossp ships with standard Postgres installs.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS todos (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    -- CHECK constraint enforces valid values without a rigid Postgres ENUM type.
    -- Adding a new status later only requires updating the CHECK, not an ALTER TYPE.
    status      VARCHAR(20)  NOT NULL DEFAULT 'PENDING'
                    CHECK (status IN ('PENDING', 'IN_PROGRESS', 'DONE')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
