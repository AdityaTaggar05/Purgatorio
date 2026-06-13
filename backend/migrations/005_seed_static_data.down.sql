DELETE FROM troops WHERE id IN ('stone-bearer', 'ashwalker', 'hoarder', 'ravager', 'coveter');
DELETE FROM building_levels WHERE building_id IN ('bastion', 'angel-spire', 'lament-basin', 'sanctum', 'barracks');
DELETE FROM building_limits WHERE building_id IN ('bastion', 'angel-spire', 'lament-basin', 'sanctum', 'barracks');
DELETE FROM buildings WHERE id IN ('bastion', 'angel-spire', 'lament-basin', 'sanctum', 'barracks');
