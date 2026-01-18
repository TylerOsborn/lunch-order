ALTER TABLE users
ADD COLUMN email_hash VARCHAR(64),
ADD COLUMN email_encrypted TEXT,
ADD COLUMN google_id_hash VARCHAR(64),
ADD COLUMN google_id_encrypted TEXT;

CREATE INDEX idx_users_email_hash ON users(email_hash);
CREATE INDEX idx_users_google_id_hash ON users(google_id_hash);

-- We are dropping the old columns. 
-- WARNING: This will delete existing data in these columns.
ALTER TABLE users
DROP COLUMN email,
DROP COLUMN google_id;
