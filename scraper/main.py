import undetected_chromedriver as uc
import requests
from backend import goodreads, librarything, olshim, data
driver = uc.Chrome(use_subprocess=False, version_main=140)


def main():
    print("Getting reviews")

    # librarything.LTScraper().getReviews("9780375822070", driver)
    # print(f"# of reviews: {len(revs)}")
    # print("posting reviews")

    # print(olshim.GenerateWorkData("OL45804W").author)
    ScrapeAndPost("9780375822070")


def ScrapeAndPost(isbn: str):
    scrapers = [goodreads.GRScraper, librarything.LTScraper]

    # get metadata, and initialize structs
    isbn, olid = olshim.ISBNtoOLIDW(isbn)
    route = data.ISBNRouteData(isbn, olid)
    work = olshim.GenerateWorkData(olid)
    reviews_final = []

    # actually scrape data
    for s in scrapers:
        reviews_final = reviews_final+s.getReviews(isbn, driver)
    # fix olid and post REVIEWS to the db
    for review in reviews_final:
        review.olid = olid
        requests.post(
            "http://127.0.0.1:1024/api/insert/reviewsingle", json=review.asJSON())
    # TODO: post ROUTING


main()
