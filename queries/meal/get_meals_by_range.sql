SELECT * FROM meals 
WHERE date >= ? AND date <= ? AND deleted_at IS NULL;