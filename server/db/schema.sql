CREATE TABLE if not exists works (
    olid TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    cover TEXT,
    description TEXT
); CREATE TABLE isbns (
    isbn TEXT PRIMARY KEY,
    olid TEXT NOT NULL REFERENCES works(olid),
    UNIQUE(isbn, olid)
);

CREATE TABLE if not exists reviews (
    review_id INTEGER PRIMARY KEY AUTOINCREMENT,
    olid TEXT NOT NULL REFERENCES works(olid),
    source TEXT NOT NULL,
    external_id TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    rating REAL NOT NULL,
    text TEXT,
    UNIQUE(external_id, text, source, username)
);

CREATE TABLE if not exists stats (
    olid TEXT PRIMARY KEY references works(olid),
    review_count INTEGER NOT NULL,
    avg_rating REAL NOT NULL,
    med_rating REAL NOT NULL,
    ci_bound REAL NOT NULL,
    UNIQUE(olid)
);

CREATE TABLE if not exists genres (
    genre_id INTEGER PRIMARY KEY AUTOINCREMENT,
    genre_name VARCHAR(20) UNIQUE NOT NULL,
    UNIQUE(genre_name)
);

CREATE TABLE bookgenres (
    olid TEXT NOT NULL REFERENCES works(olid),
    genre_id INTEGER NOT NULL REFERENCES genres(genre_id),
    UNIQUE(olid, genre_id)
);

CREATE VIEW overall_rating as
SELECT r.olid,AVG(r.positive) FROM REVIEWS r GROUP BY r.olid;
