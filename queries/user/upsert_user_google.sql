INSERT INTO users (name, email_hash, email_encrypted, google_id_hash, google_id_encrypted, first_name, last_name, avatar_url, is_admin)
VALUES (:name, :email_hash, :email_encrypted, :google_id_hash, :google_id_encrypted, :first_name, :last_name, :avatar_url, :is_admin)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    email_hash = VALUES(email_hash),
    email_encrypted = VALUES(email_encrypted),
    first_name = VALUES(first_name),
    last_name = VALUES(last_name),
    avatar_url = VALUES(avatar_url),
    is_admin = VALUES(is_admin),
    updated_at = CURRENT_TIMESTAMP;