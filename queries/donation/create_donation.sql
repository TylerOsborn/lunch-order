INSERT INTO donations (created_at, updated_at, meal_id, donor_id)
SELECT NOW(), NOW(), ?, ?
    WHERE NOT EXISTS (
        SELECT 1 FROM donations WHERE donor_id = ? AND DATE(created_at) = CURDATE()
    );