CREATE TABLE user (
  id INTEGER PRIMARY KEY NOT NULL,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  approved boolean NOT NULL DEFAULT false,
  created_at INTEGER DEFAULT (strftime('%s', 'now'))
);
