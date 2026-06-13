ALTER TABLE battles ALTER COLUMN outcome DROP DEFAULT;
ALTER TABLE battles ALTER COLUMN outcome TYPE TEXT USING outcome::TEXT;
ALTER TABLE battles ALTER COLUMN outcome SET DEFAULT 'pending';

ALTER TABLE battles ADD CONSTRAINT battles_outcome_check CHECK (outcome IN ('victory', 'defeat', 'threshold_failed', 'pending'));

ALTER TABLE battles ALTER COLUMN destruction TYPE FLOAT USING destruction::FLOAT;
ALTER TABLE battles ADD COLUMN IF NOT EXISTS finished_at TIMESTAMPTZ;
ALTER TABLE battles ALTER COLUMN base_snapshot_id DROP NOT NULL;

ALTER TABLE base_snapshots RENAME COLUMN created_at TO captured_at;

DROP TYPE IF EXISTS battle_outcome;
