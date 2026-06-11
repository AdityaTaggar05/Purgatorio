ALTER TABLE base_layouts ADD level INT NOT NULL DEFAULT 1 CHECK (level >= 1);
