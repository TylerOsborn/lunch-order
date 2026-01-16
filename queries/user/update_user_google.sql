UPDATE users
SET name = :name,
    email = :email,
    first_name = :first_name,
    last_name = :last_name,
    avatar_url = :avatar_url,
    is_admin = :is_admin,
    updated_at = CURRENT_TIMESTAMP
WHERE google_id = :google_id;
