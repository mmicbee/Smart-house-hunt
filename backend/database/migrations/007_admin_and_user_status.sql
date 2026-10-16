ALTER TABLE users ADD COLUMN status TEXT NOT NULL DEFAULT 'active';

INSERT OR IGNORE INTO users (role_id, name, email, password, status)
VALUES (
    4,
    'System Administrator',
    'admin@smarthousehunt.com',
    '$2a$10$LjG.J/E9cY4LXqA5U4crIeeM2uNNnRFJZBTssbJ1aiAsVNsLl8cjy',
    'active'
);
