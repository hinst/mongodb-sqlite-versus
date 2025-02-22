PRAGMA journal_mode = WAL;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    email TEXT NOT NULL,
    createdAt TEXT NOT NULL,
    level INTEGER NOT NULL
);

VACUUM;