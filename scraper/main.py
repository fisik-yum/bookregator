import undetected_chromedriver as uc
import requests
import logging, coloredlogs
from backend import goodreads, librarything, olshim, data
driver = uc.Chrome(use_subprocess=False, version_main=139)
logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
coloredlogs.install()
logger=logging.getLogger(__name__)

def main():
    # librarything.LTScraper().getReviews("9780375822070", driver)
    # print(f"# of reviews: {len(revs)}")
    # print("posting reviews")

    # print(olshim.GenerateWorkData("OL45804W").author)
    books=["9780375822070"," 9780765348784","9781250887672"]
    for book in books:
        ScrapeAndPost(book)


def ScrapeAndPost(isbn: str):
    # get list of scrapers
    scrapers = [goodreads.GRScraper, librarything.LTScraper]

    # get metadata, and initialize structs
    isbn, olid = olshim.ISBNtoOLIDW(isbn)
    logger.info(f"Working on (ISBN,OLID), ({isbn},{olid})")

    route = data.ISBNRouteData(isbn, olid)
    work = olshim.GenerateWorkData(olid)
    reviews_final = []

    # actually scrape data
    for s in scrapers:
        reviews_final = reviews_final+s.getReviews(isbn, driver)
    # insert work
    requests.post(
        "http://127.0.0.1:1024/api/insert/work", json=work.asJSON())
    # create routing
    requests.post(
        "http://127.0.0.1:1024/api/insert/route", json=route.asJSON())
    # fix olid and post REVIEWS to the db
    for review in reviews_final:
        review.olid = olid
        requests.post(
            "http://127.0.0.1:1024/api/insert/reviewsingle", json=review.asJSON())


main()
