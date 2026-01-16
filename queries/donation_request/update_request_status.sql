UPDATE donation_requests 
SET status = ?, donation_id = ?, updated_at = NOW() 
WHERE id = ? AND deleted_at IS NULL;