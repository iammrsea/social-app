-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('ADMIN', 'REGULAR', 'MODERATOR', 'GUEST')),
    reputation_score INT NOT NULL DEFAULT 0,
    badges TEXT[], -- Array of strings for badges
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    banned_at TIMESTAMP,
    ban_start_date TIMESTAMP,
    ban_end_date TIMESTAMP,
    reason_for_ban TEXT,
    is_ban_indefinite BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Optional: Seed initial data
INSERT INTO users (id, username, email, role, reputation_score, badges, is_banned, created_at, updated_at)
VALUES
    ('cuid1', 'johndoe', 'johndoe@example.com', 'REGULAR', 100, ARRAY['badge1', 'badge2'], FALSE, NOW(), NOW()),
    ('cuid2', 'janedoe', 'janedoe@example.com', 'MODERATOR', 200, ARRAY['badge3'], FALSE, NOW(), NOW())
ON CONFLICT DO NOTHING;