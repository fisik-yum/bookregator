import sys

import joblib
import pandas as pd
from sklearn.compose import ColumnTransformer
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics import classification_report, f1_score
from sklearn.model_selection import GridSearchCV, train_test_split
from sklearn.pipeline import Pipeline
from sklearn.svm import LinearSVC

import db

DEFAULT_OUT_PATH = "reviewclassifier.joblib"


def load_training_data(labels_csv: str, db_path: str) -> pd.DataFrame:
    labels = pd.read_csv(labels_csv)
    con = db.connect(db_path)
    reviews = pd.read_sql_query(
        "SELECT review_id, rating, text FROM reviews WHERE text IS NOT NULL", con
    )
    merged = reviews.merge(labels[["review_id", "label"]], on="review_id")
    merged["label"] = merged["label"].map({"POSITIVE": 1, "NEGATIVE": 0})
    return merged.dropna(subset=["text", "rating", "label"])


def train(df: pd.DataFrame):
    X_train, X_test, y_train, y_test = train_test_split(
        df[["text", "rating"]],
        df["label"],
        test_size=0.25,
        stratify=df["label"],
    )

    preprocessor = ColumnTransformer(
        transformers=[
            ("text_tfidf", TfidfVectorizer(stop_words="english"), "text"),
        ]
    )
    pipeline = Pipeline(
        [
            ("preprocessor", preprocessor),
            ("svc", LinearSVC()),
        ]
    )

    param_grid = {
        "preprocessor__text_tfidf__max_df": [0.8, 1.0],
        "preprocessor__text_tfidf__ngram_range": [(1, 1), (1, 2)],
        "svc__C": [1.0, 1.2, 1.5, 1.9],
    }

    grid_search = GridSearchCV(
        pipeline, param_grid, cv=5, scoring="f1", n_jobs=-1, verbose=1
    )

    print("Starting Grid Search CV...")
    grid_search.fit(X_train, y_train)

    print(f"Best CV F1-Score: {grid_search.best_score_:.4f}")
    print(f"Best Parameters: {grid_search.best_params_}")

    best_model = grid_search.best_estimator_
    y_pred = best_model.predict(X_test)
    print(f"Test F1-Score: {f1_score(y_test, y_pred):.4f}")
    print(classification_report(y_test, y_pred))
    return best_model


def main():
    if len(sys.argv) < 2:
        print("Usage: python train_sentiment.py <bert_labels.csv> [db_path] [out_joblib]")
        sys.exit(1)

    labels_csv = sys.argv[1]
    db_path = sys.argv[2] if len(sys.argv) > 2 else db.DEFAULT_DB_PATH
    out_path = sys.argv[3] if len(sys.argv) > 3 else DEFAULT_OUT_PATH

    training_data = load_training_data(labels_csv, db_path)
    model = train(training_data)
    joblib.dump(model, out_path)
    print(f"Saved model to {out_path}")


if __name__ == "__main__":
    main()
