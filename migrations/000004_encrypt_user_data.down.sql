ALTER TABLE users
ADD COLUMN email VARCHAR(255) UNIQUE,
ADD COLUMN google_id VARCHAR(255) UNIQUE;

ALTER TABLE users
DROP COLUMN email_hash,
DROP COLUMN email_encrypted,
DROP COLUMN google_id_hash,
DROP COLUMN google_id_encrypted;
