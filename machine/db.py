import os
import sqlite3

DEFAULT_DB_PATH = os.path.join(
    os.path.dirname(__file__), "..", "server", "artifacts", "test.sqlite3"
)

REVIEW_SENTIMENT_SCHEMA = """
CREATE TABLE IF NOT EXISTS review_sentiment (
    review_id INTEGER PRIMARY KEY REFERENCES reviews(review_id),
    label TEXT NOT NULL,
    score REAL NOT NULL
)
"""

BOOK_RECOMMENDATIONS_SCHEMA = """
CREATE TABLE IF NOT EXISTS book_recommendations (
    olid TEXT NOT NULL REFERENCES works(olid),
    similar_olid TEXT NOT NULL REFERENCES works(olid),
    rank INTEGER NOT NULL,
    score REAL NOT NULL,
    UNIQUE(olid, rank)
)
"""


def connect(db_path: str = DEFAULT_DB_PATH) -> sqlite3.Connection:
    con = sqlite3.connect(db_path)
    con.execute(REVIEW_SENTIMENT_SCHEMA)
    con.execute(BOOK_RECOMMENDATIONS_SCHEMA)
    con.commit()
    return con
