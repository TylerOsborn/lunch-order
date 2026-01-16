SELECT 
    d.id, 
    d.created_at, 
    d.updated_at, 
    d.meal_id, 
    d.donor_id, 
    d.recipient_id,
    m.id AS "meal.id",
    m.description AS "meal.description",
    m.date AS "meal.date",
    u.id AS "donor.id",
    u.name AS "donor.name"
FROM donations d
JOIN meals m ON d.meal_id = m.id
JOIN users u ON d.donor_id = u.id
WHERE (d.recipient_id <= 0 OR d.recipient_id IS NULL) 
AND DATE(m.date) = DATE(?);