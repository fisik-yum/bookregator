CREATE TABLE works (
  olid TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT,
  description TEXT,
  published_year INTEGER
);

CREATE TABLE isbns (
  isbn TEXT PRIMARY KEY,
  olid TEXT NOT NULL REFERENCES works(olid)
);

CREATE TABLE reviews (
  review_id INTEGER PRIMARY KEY AUTOINCREMENT,
  olid TEXT NOT NULL REFERENCES works(olid),
  source TEXT NOT NULL,
  external_id TEXT NOT NULL,
  username TEXT NOT NULL,
  rating REAL,
  text TEXT,
  UNIQUE(olid, source, external_id, text)
);
