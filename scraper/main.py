import zendriver
import itertools
#import undetected_chromedriver as uc
import requests
import logging
import coloredlogs
import zendriver as z
from backend import goodreads, librarything, olshim, data
import asyncio
#driver = uc.Chrome(use_subprocess=False, version_main=139)
logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
coloredlogs.install()
logger = logging.getLogger(__name__)


def main():
    olids=set()
    isbns=set()
    with open("isbns.txt") as f:
        lc=0
        for l in f:
            print(f"Working on line {lc}")
            book=l.strip(" \n\r")
            try:
                ScrapeAndPost(book,olids,isbns,browser)
            except:
                lc+=1
                continue
            lc+=1


async def ScrapeAndPost(isbn: str,olids: set, isbns: set):
    # auth setup
    authkeys = ("scraper", "opensesame")
    # get list of scrapers
    scrapers = [goodreads.GRScraper]
    if isbn in isbns:
        return

    isbns.add(isbn)
    # get metadata, and initialize structs
    isbn, olid = olshim.ISBNtoOLIDW(isbn)

    if olid not in olids:
        return

    logger.info(f"Working on (ISBN,OLID), ({isbn},{olid})")
    olids.add(olid)


    route = data.ISBNRouteData(isbn, olid)
    work = olshim.GenerateWorkData(olid)

    browser=await z.start()
    tasks=[]
    # actually scrape data
    for s in scrapers:
        tab=await browser.get("",new_tab=True)
        asyncs.append(s.getReviews(isbn, tab))
    # fix olid
    for review in reviews_final:
        review.olid = olid
    # insert work
    reviews_final=asyncio.gather(*tasks)
    reviews_final=list(itertools.chain.from_iterable(reviews_final))
    print(reviews_final)

    requests.post(
        "http://127.0.0.1:1024/api/insert/work", json=work.asJSON(), auth=authkeys)
    # create routing
    requests.post(
        "http://127.0.0.1:1024/api/insert/route", json=route.asJSON(), auth=authkeys)
    # Post REVIEWS to the db using reviewmultiple for decreased overhead
    requests.post(
        "http://127.0.0.1:1024/api/insert/reviewmultiple", json=[r.asJSON() for r in reviews_final], auth=authkeys)
main()
