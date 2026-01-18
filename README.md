# Lunch Order App

A Go-based lunch ordering system.

## Database Security

This application uses **Application-Level Encryption** (Blind Indexing) to protect sensitive user data (`email` and `google_id`).

*   **Storage:** Data is encrypted using AES-GCM before being saved to the database.
*   **Search:** A deterministic HMAC-SHA256 hash is stored in a separate column (`_hash`) to allow for efficient lookups without exposing the raw data.

### Configuration

You must set the `DATA_ENCRYPTION_KEY` environment variable. This key is used for both encryption and hashing.

**Generate a new key:**
```bash
openssl rand -hex 32
```

Add this to your `.env` file:
```bash
DATA_ENCRYPTION_KEY=your_generated_64_char_hex_string
```

### How to Query the Database

Since data is encrypted, you cannot simply run `SELECT * FROM users WHERE email = 'user@example.com'`. You must use the helper tool to generate the hash or decrypt the data.

#### 1. Finding a User (Generating a Hash)
To write a SQL query for a specific user, you first need to generate the hash of their email or Google ID.

```bash
# Generate hash for an email
go run tools/crypto_tool.go -action=hash -input="tyler@example.com"
```

**Output:**
```text
Hash: 8f4b2e...
```

**SQL Query:**
```sql
SELECT * FROM users WHERE email_hash = '8f4b2e...';
```

#### 2. Reading Encrypted Data
If you query the database and see an encrypted string (e.g., in the `email_encrypted` column), you can decrypt it to see the original value.

```bash
# Decrypt a value
go run tools/crypto_tool.go -action=decrypt -input="<encrypted_string_from_db>"
```

**Output:**
```text
Decrypted: tyler@example.com
```

#### 3. Encrypting Data Manually
If you need to manually insert a user via SQL, you can generate the encrypted value.

```bash
go run tools/crypto_tool.go -action=encrypt -input="tyler@example.com"
```
