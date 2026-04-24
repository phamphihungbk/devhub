ALTER TABLE deployments
    DROP COLUMN IF EXISTS runner_error,
    DROP COLUMN IF EXISTS runner_output;
