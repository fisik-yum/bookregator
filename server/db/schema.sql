CREATE TABLE works (
    olid TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    cover TEXT,
    description TEXT
);

CREATE TABLE isbns (
    isbn TEXT PRIMARY KEY,
    olid TEXT NOT NULL REFERENCES works(olid),
    UNIQUE(isbn, olid)
);

CREATE TABLE reviews (
    review_id INTEGER PRIMARY KEY AUTOINCREMENT,
    olid TEXT NOT NULL REFERENCES works(olid),
    source TEXT NOT NULL,
    external_id TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    rating REAL,
    text TEXT,
    UNIQUE(external_id, text, source, username)
);

CREATE TABLE stats (
    olid TEXT PRIMARY KEY references works(olid),
    review_count INTEGER,
    rating REAL,
    ci_bound REAL,
    UNIQUE(olid)
);
