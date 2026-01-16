SELECT * FROM users 
WHERE name = ? AND deleted_at IS NULL 
LIMIT 1;