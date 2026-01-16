INSERT INTO donation_requests (created_at, updated_at, requester_id, status) 
VALUES (NOW(), NOW(), ?, ?);