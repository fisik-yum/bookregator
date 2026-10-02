import sys

import joblib
import pandas as pd

import db

DEFAULT_MODEL_PATH = "reviewclassifier.joblib"


def main():
    model_path = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_MODEL_PATH
    db_path = sys.argv[2] if len(sys.argv) > 2 else db.DEFAULT_DB_PATH

    model = joblib.load(model_path)
    con = db.connect(db_path)
    reviews = pd.read_sql_query(
        "SELECT review_id, rating, text FROM reviews WHERE text IS NOT NULL", con
    )

    predictions = model.decision_function(reviews[["text", "rating"]])
    reviews["label"] = ["POSITIVE" if p > 0 else "NEGATIVE" for p in predictions]
    reviews["score"] = predictions

    con.executemany(
        """
        INSERT INTO review_sentiment (review_id, label, score)
        VALUES (?, ?, ?)
        ON CONFLICT(review_id) DO UPDATE SET label=excluded.label, score=excluded.score
        """,
        reviews[["review_id", "label", "score"]].itertuples(index=False, name=None),
    )
    con.commit()
    print(f"Tagged {len(reviews)} reviews")


if __name__ == "__main__":
    main()
