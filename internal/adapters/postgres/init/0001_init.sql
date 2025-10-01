CREATE TABLE IF NOT EXISTS meetings (
  id           UUID PRIMARY KEY,
  title        VARCHAR(200)      NOT NULL,
  starts_at    TIMESTAMPTZ       NOT NULL,
  duration_sec INTEGER           NOT NULL,            
  status       TEXT              NOT NULL,          
  created_at   TIMESTAMPTZ       NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ       NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_meetings_starts_at ON meetings(starts_at);

DROP TABLE IF EXISTS meetings;