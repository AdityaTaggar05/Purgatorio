CREATE TYPE battle_outcome AS ENUM ('victory', 'defeat', 'threshold_failed');

ALTER TABLE battles ALTER COLUMN outcome DROP DEFAULT;
ALTER TABLE battles DROP CONSTRAINT IF EXISTS battles_outcome_check;
ALTER TABLE battles ALTER COLUMN outcome TYPE battle_outcome USING outcome::battle_outcome;

ALTER TABLE battles ALTER COLUMN destruction TYPE INT USING destruction::INT;
ALTER TABLE battles ALTER COLUMN destruction DROP DEFAULT;
ALTER TABLE battles ALTER COLUMN loot DROP DEFAULT;
ALTER TABLE battles ALTER COLUMN duration DROP DEFAULT;
ALTER TABLE battles DROP COLUMN IF EXISTS finished_at;
ALTER TABLE battles ALTER COLUMN base_snapshot_id SET NOT NULL;

ALTER TABLE base_snapshots RENAME COLUMN captured_at TO created_at;
