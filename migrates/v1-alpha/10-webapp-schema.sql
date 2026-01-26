SET search_path TO webapp;

-- Add unique constraint to studies.user_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'studies_user_id_unique_constraint'
    ) THEN
        ALTER TABLE studies
        ADD CONSTRAINT studies_user_id_unique_constraint UNIQUE (user_id);
    END IF;
END$$;

-- Make degree_plan_code nullable if it is currently NOT NULL
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'webapp'
          AND table_name = 'studies'
          AND column_name = 'degree_plan_code'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE studies
        ALTER degree_plan_code DROP NOT NULL;
    END IF;
END$$;

-- Drop start_year column if it exists
ALTER TABLE studies DROP COLUMN IF EXISTS start_year;