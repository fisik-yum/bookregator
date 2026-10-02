import sys

import pandas as pd
from sklearn.neighbors import NearestNeighbors
from sklearn.preprocessing import StandardScaler

import db

DEFAULT_N_NEIGHBORS = 10


def load_features(con) -> pd.DataFrame:
    stats = pd.read_sql_query(
        "SELECT olid, review_count, avg_rating, med_rating FROM stats", con
    )
    sentiment = pd.read_sql_query(
        """
        SELECT r.olid, AVG(CASE WHEN rs.label = 'POSITIVE' THEN 1.0 ELSE 0.0 END) AS positive_ratio
        FROM reviews r
        JOIN review_sentiment rs ON rs.review_id = r.review_id
        GROUP BY r.olid
        """,
        con,
    )
    features = stats.merge(sentiment, on="olid", how="left")
    features["positive_ratio"] = features["positive_ratio"].fillna(0.5)
    return features


def build_neighbors(features: pd.DataFrame, n_neighbors: int) -> pd.DataFrame:
    feature_cols = ["review_count", "avg_rating", "med_rating", "positive_ratio"]
    X = StandardScaler().fit_transform(features[feature_cols])

    k = min(n_neighbors + 1, len(features))
    nn = NearestNeighbors(n_neighbors=k)
    nn.fit(X)
    distances, indices = nn.kneighbors(X)

    rows = []
    for row_idx, (dist_row, idx_row) in enumerate(zip(distances, indices)):
        olid = features.iloc[row_idx]["olid"]
        rank = 0
        for dist, neighbor_idx in zip(dist_row, idx_row):
            if neighbor_idx == row_idx:
                continue
            rank += 1
            rows.append(
                {
                    "olid": olid,
                    "similar_olid": features.iloc[neighbor_idx]["olid"],
                    "rank": rank,
                    "score": float(dist),
                }
            )
    return pd.DataFrame(rows)


def main():
    db_path = sys.argv[1] if len(sys.argv) > 1 else db.DEFAULT_DB_PATH
    n_neighbors = int(sys.argv[2]) if len(sys.argv) > 2 else DEFAULT_N_NEIGHBORS

    con = db.connect(db_path)
    features = load_features(con)
    recommendations = build_neighbors(features, n_neighbors)

    con.execute("DELETE FROM book_recommendations")
    con.executemany(
        "INSERT INTO book_recommendations (olid, similar_olid, rank, score) VALUES (?, ?, ?, ?)",
        recommendations[["olid", "similar_olid", "rank", "score"]].itertuples(
            index=False, name=None
        ),
    )
    con.commit()
    print(f"Wrote {len(recommendations)} recommendations for {len(features)} books")


if __name__ == "__main__":
    main()
