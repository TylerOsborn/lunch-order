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
    donor.id AS "donor.id",
    donor.name AS "donor.name",
    recipient.id AS "recipient.id",
    recipient.name AS "recipient.name"
FROM donations d
JOIN meals m ON d.meal_id = m.id
JOIN users donor ON d.donor_id = donor.id
LEFT JOIN users recipient ON d.recipient_id = recipient.id
WHERE DATE(d.created_at) = DATE(?);