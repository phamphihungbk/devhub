ALTER TABLE deployments
    ADD COLUMN IF NOT EXISTS runner_output TEXT,
    ADD COLUMN IF NOT EXISTS runner_error TEXT;
