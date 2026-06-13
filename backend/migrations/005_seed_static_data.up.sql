INSERT INTO buildings (id, name, size, price, currency, category) VALUES
  ('bastion', 'Bastion', 1, 50, 'penitence', 'defense'),
  ('angel-spire', 'Angel Spire', 2, 500, 'penitence', 'defense'),
  ('lament-basin', 'Lament Basin', 2, 300, 'penitence', 'resource'),
  ('sanctum', 'Sanctum', 3, 10000, 'penitence', 'other'),
  ('barracks', 'Barracks', 3, 500, 'penitence', 'army')
ON CONFLICT (id) DO NOTHING;

INSERT INTO building_limits (building_id, terrace_level, max_allowed) VALUES
  ('bastion', 1, 10), ('bastion', 2, 20), ('bastion', 3, 35),
  ('angel-spire', 1, 1), ('angel-spire', 2, 2), ('angel-spire', 3, 3),
  ('lament-basin', 1, 1), ('lament-basin', 2, 2), ('lament-basin', 3, 3),
  ('sanctum', 1, 1), ('sanctum', 2, 1), ('sanctum', 3, 1),
  ('barracks', 1, 1), ('barracks', 2, 1), ('barracks', 3, 1)
ON CONFLICT (building_id, terrace_level) DO NOTHING;

INSERT INTO building_levels (building_id, level, hp, damage_per_second, production_rate, storage_capacity, attack_range, upgrade_cost, upgrade_time) VALUES
  ('bastion', 1, 300, NULL, NULL, NULL, 0, 50, 60),
  ('bastion', 2, 500, NULL, NULL, NULL, 0, 150, 300),
  ('bastion', 3, 800, NULL, NULL, NULL, 0, 400, 900),
  ('bastion', 4, 1200, NULL, NULL, NULL, 0, 900, 1800),
  ('angel-spire', 1, 450, 12, NULL, NULL, 5.0, 500, 600),
  ('angel-spire', 2, 600, 18, NULL, NULL, 5.0, 1200, 1800),
  ('angel-spire', 3, 800, 26, NULL, NULL, 5.0, 2500, 3600),
  ('angel-spire', 4, 1050, 36, NULL, NULL, 5.0, 5000, 7200),
  ('angel-spire', 5, 1400, 50, NULL, NULL, 5.0, 9000, 14400),
  ('lament-basin', 1, 400, NULL, 2, 500, 0, 300, 600),
  ('lament-basin', 2, 550, NULL, 4, 1000, 0, 800, 1800),
  ('lament-basin', 3, 700, NULL, 6, 1800, 0, 2000, 3600),
  ('lament-basin', 4, 900, NULL, 8, 3000, 0, 4500, 7200),
  ('lament-basin', 5, 1150, NULL, 12, 5000, 0, 9000, 14400),
  ('sanctum', 1, 500, NULL, NULL, NULL, 0, 5000, 3600),
  ('sanctum', 2, 600, NULL, NULL, NULL, 0, 12000, 7200),
  ('sanctum', 3, 700, NULL, NULL, NULL, 0, 25000, 14400),
  ('sanctum', 4, 800, NULL, NULL, NULL, 0, 50000, 28800),
  ('sanctum', 5, 900, NULL, NULL, NULL, 0, 100000, 57600),
  ('barracks', 1, 400, NULL, NULL, 50, 0, 500, 600),
  ('barracks', 2, 550, NULL, NULL, 100, 0, 1500, 1800),
  ('barracks', 3, 700, NULL, NULL, 200, 0, 4000, 3600)
ON CONFLICT (building_id, level) DO NOTHING;

INSERT INTO troops (id, name, training_cost, space, hp, dps, speed, attack_range, preferred_target) VALUES
  ('stone-bearer', 'Stone Bearer', 80, 6, 300, 12, 1.0, 1.0, 'defense'),
  ('ashwalker', 'Ashwalker', 60, 3, 80, 25, 2.0, 2.0, 'defense'),
  ('hoarder', 'Hoarder', 50, 4, 120, 18, 1.5, 1.5, 'resource'),
  ('ravager', 'Ravager', 90, 7, 200, 40, 1.2, 1.0, 'defense'),
  ('coveter', 'Coveter', 30, 2, 60, 15, 2.5, 3.0, 'any')
ON CONFLICT (id) DO NOTHING;
