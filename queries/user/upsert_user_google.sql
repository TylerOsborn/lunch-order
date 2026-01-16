INSERT INTO users (name, email, google_id, first_name, last_name, avatar_url, is_admin)
VALUES (:name, :email, :google_id, :first_name, :last_name, :avatar_url, :is_admin)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    email = VALUES(email),
    first_name = VALUES(first_name),
    last_name = VALUES(last_name),
    avatar_url = VALUES(avatar_url),
    is_admin = VALUES(is_admin),
    updated_at = CURRENT_TIMESTAMP;
