# bert_tag.py and train_sentiment.py are run manually, not part of this pipeline.

import build_recommendations
import tag_reviews


def main():
    tag_reviews.main()
    build_recommendations.main()


if __name__ == "__main__":
    main()
