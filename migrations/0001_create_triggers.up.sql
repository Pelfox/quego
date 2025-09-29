CREATE TABLE IF NOT EXISTS triggers (
  id BLOB(16) PRIMARY KEY,
  function_name TEXT NOT NULL,
  trigger_type NOT NULL CHECK (trigger_type in ('EVENT', 'CRON')),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_triggers_function_name ON triggers(function_name);
