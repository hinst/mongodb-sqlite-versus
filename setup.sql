CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    accessToken TEXT NOT NULL,
    email TEXT NOT NULL,
    createdAt INTEGER NOT NULL,
    level INTEGER NOT NULL
);
