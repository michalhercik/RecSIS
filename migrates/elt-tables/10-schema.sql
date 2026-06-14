DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_namespace WHERE nspname = 'recommender'
  ) THEN
    EXECUTE 'CREATE SCHEMA recommender AUTHORIZATION recommender';
  END IF;
END
$$;
