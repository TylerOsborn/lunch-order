SELECT 
    dr.id, 
    dr.created_at, 
    dr.updated_at, 
    dr.deleted_at, 
    dr.requester_id, 
    dr.status, 
    dr.donation_id,
    u.id AS "requester.id",
    u.name AS "requester.name",
    d.id AS "donation.id",
    d.meal_id AS "donation.meal_id",
    d.donor_id AS "donation.donor_id",
    donor.id AS "donation.donor.id",
    donor.name AS "donation.donor.name",
    m.id AS "donation.meal.id",
    m.description AS "donation.meal.description",
    m.date AS "donation.meal.date"
FROM donation_requests dr
JOIN users u ON dr.requester_id = u.id
LEFT JOIN donations d ON dr.donation_id = d.id
LEFT JOIN users donor ON d.donor_id = donor.id
LEFT JOIN meals m ON d.meal_id = m.id
WHERE dr.requester_id = ? 
AND dr.status = 'pending' 
AND DATE(dr.created_at) = DATE(?)
AND dr.deleted_at IS NULL
ORDER BY dr.created_at DESC;