UPDATE donations 
SET recipient_id = ?, updated_at = NOW() 
WHERE id = ? AND (recipient_id = 0 OR recipient_id IS NULL);