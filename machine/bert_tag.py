import sys

import pandas as pd
from transformers import pipeline

import db


def main():
    db_path = sys.argv[1] if len(sys.argv) > 1 else db.DEFAULT_DB_PATH
    out_csv = sys.argv[2] if len(sys.argv) > 2 else "reviews_tagged.csv"

    con = db.connect(db_path)
    reviews = pd.read_sql_query(
        "SELECT review_id, text FROM reviews WHERE text IS NOT NULL", con
    )

    sentiment_pipeline = pipeline("sentiment-analysis", truncation=True)
    results = sentiment_pipeline(reviews["text"].fillna("").tolist())
    reviews["label"] = [r["label"] for r in results]
    reviews["score"] = [r["score"] for r in results]

    reviews[["review_id", "label", "score"]].to_csv(out_csv, index=False)
    print(f"Wrote {len(reviews)} labeled reviews to {out_csv}")


if __name__ == "__main__":
    main()
