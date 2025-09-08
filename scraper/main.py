import undetected_chromedriver as uc
import requests
from backend import goodreads, librarything
driver = uc.Chrome(use_subprocess=False, version_main=138)


def main():
    print("Getting reviews")
    librarything.LTScraper().ISBNtoLTCode("9780375822070", driver)
    # print(f"# of reviews: {len(revs)}")
    # print("posting reviews")
    # for review in revs:
    #    review.olid = "OL45804"
    #    print(review.text)
    #    requests.post(
    #        "http://127.0.0.1:1024/api/insert/reviewsingle", json=review.asJSON())


main()
