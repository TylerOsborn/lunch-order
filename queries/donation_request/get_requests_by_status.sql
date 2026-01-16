SELECT 
    dr.id, 
    dr.created_at, 
    dr.updated_at, 
    dr.deleted_at, 
    dr.requester_id, 
    dr.status, 
    dr.donation_id,
    u.id AS "requester.id",
    u.name AS "requester.name"
FROM donation_requests dr
JOIN users u ON dr.requester_id = u.id
WHERE dr.status = ? 
AND dr.deleted_at IS NULL
ORDER BY dr.created_at ASC;