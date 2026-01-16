SELECT * FROM meals 
WHERE description = ? AND date = ? AND deleted_at IS NULL 
LIMIT 1;