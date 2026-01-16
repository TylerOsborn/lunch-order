SELECT m.* 
FROM meals m
JOIN donation_request_meals drm ON m.id = drm.meal_id
WHERE drm.donation_request_id = ?;