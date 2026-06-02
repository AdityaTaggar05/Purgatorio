CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE item_category AS ENUM ('defense', 'army', 'resource', 'other');
CREATE TYPE currency_type AS ENUM ('penitence', 'grace');
CREATE TYPE battle_outcome AS ENUM ('victory', 'defeat', 'threshold_failed');

-- Auth and User Management
CREATE TABLE auth(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_on TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users(
  id UUID PRIMARY KEY REFERENCES auth(id) ON DELETE CASCADE,
  username TEXT UNIQUE NOT NULL CHECK (length(username) BETWEEN 3 and 14),
  xp INT NOT NULL DEFAULT 0 CHECK (xp >= 0),
  level INT GENERATED ALWAYS AS (
    GREATEST(1, FLOOR(SQRT(xp / 25.0)))::INT
  ) STORED,
  terrace_level INT NOT NULL DEFAULT 1 CHECK (terrace_level >= 1),
  created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_stats(
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  attacks INT NOT NULL DEFAULT 0 CHECK (attacks >= 0),
  attacks_success INT NOT NULL DEFAULT 0 CHECK (attacks_success >= 0),
  defenses INT NOT NULL DEFAULT 0 CHECK (defenses >= 0),
  defenses_success INT NOT NULL DEFAULT 0 CHECK (defenses_success >= 0),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_economy(
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  penitence INT NOT NULL DEFAULT 500 CHECK (penitence >= 0),
  grace INT NOT NULL DEFAULT 50 CHECK (grace >= 0),
  max_penitence INT NOT NULL DEFAULT 5000 CHECK (max_penitence > 0),
  collector_pending_penitence INT NOT NULL DEFAULT 0 CHECK (collector_pending_penitence >= 0),
  collector_reset_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_combat (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  sin_meter INT NOT NULL DEFAULT 0 CHECK (sin_meter BETWEEN 0 AND 100),
  last_attack_at TIMESTAMPTZ,
  shield_expires_at TIMESTAMPTZ,
  shield_max_duration INT NOT NULL DEFAULT 28800,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Shop Management
CREATE TABLE items(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT UNIQUE NOT NULL,
  price INT NOT NULL CHECK (price >= 0),
  currency currency_type NOT NULL DEFAULT 'penitence',
  category item_category NOT NULL DEFAULT 'other',
  created_on TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE item_limits(
  item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  terrace_level INT NOT NULL CHECK (terrace_level >= 1),
  max_allowed INT NOT NULL CHECK (max_allowed >= 0),
  PRIMARY KEY (item_id, terrace_level)
);

CREATE TABLE user_items(
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  item_id UUID REFERENCES items(id) ON DELETE CASCADE,
  quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  PRIMARY KEY(user_id, item_id)
);

-- State Management
CREATE TABLE game_state(
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  scene_id TEXT,
  fallback TEXT,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE base_layouts(
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  buildings JSONB NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Battle Management
CREATE TABLE base_snapshots(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- careful here, since NULL user_id would mean a deleted player. remember to add in game
  buildings JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE battles(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  attacker_id UUID REFERENCES users(id) ON DELETE SET NULL, -- same as above here
  defender_id UUID REFERENCES users(id) ON DELETE SET NULL, -- and here
  outcome battle_outcome NOT NULL,
  destruction INT NOT NULL CHECK (destruction BETWEEN 0 AND 100),
  loot INT NOT NULL CHECK (loot >= 0),
  duration INT NOT NULL CHECK (duration >= 0),
  base_snapshot_id UUID NOT NULL REFERENCES base_snapshots(id) ON DELETE CASCADE,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT same_player CHECK (attacker_id <> defender_id)
);

CREATE TABLE battle_replays(
  battle_id UUID PRIMARY KEY REFERENCES battles(id) ON DELETE CASCADE,
  data JSONB NOT NULL
);

-- Indexes
CREATE INDEX idx_battle_logs_attacker ON battles(attacker_id, started_at DESC); -- to view attack logs
CREATE INDEX idx_battle_logs_defender ON battles(defender_id, started_at DESC); -- to view defense logs
CREATE INDEX idx_player_leagues ON users(terrace_level); -- for and when leagues are implemented
