import bs4
import undetected_chromedriver as uc
import time
import logging
import coloredlogs
from . import data

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
coloredlogs.install()
global logger


class LTScraper(data.Scraper):
    url = "https://www.librarything.com/work/"
    source = "librarything"

    @staticmethod
    def ISBNtoLTCode(isbn: str, driver: uc.Chrome) -> str:
        logger = logging.getLogger(__name__)
        searchurl: str = "https://www.librarything.com/search.php?search=" + \
            isbn+"&searchtype=newwork_titles&sortchoice=0"

        driver.get(searchurl)
        time.sleep(5)
        soup = bs4.BeautifulSoup(driver.page_source, features="html.parser")
        table = soup.find(
            "table", attrs={"class": "msg", "border": "0", "cellspacing": "0"})
        if table:
            res = table.a["data-workid"].__str__()
            return res
        else:
            logger.critical("No ISBN-LibraryThing work map")
            raise Exception("No ISBN-Work map")
        return ""

    @staticmethod
    def getReviews(isbn: str, driver: uc.Chrome) -> list[data.ReviewData]:
        ret: list[data.ReviewData] = []

        logger = logging.getLogger(__name__)
        logger.debug(f"Fetching details for {isbn}")

        try:
            ltcode = LTScraper.ISBNtoLTCode(isbn, driver)
        except Exception:
            return ret

        driver.get(LTScraper.url+ltcode+"/reviews")
        time.sleep(3)
        soup = bs4.BeautifulSoup(driver.page_source, features="html.parser")

        section = soup.select_one(
            "div.mr_grid_group.mr_grid_group_type_reviews.group_1")
        if not section:
            print("no section")
        cards = section.select("div.mr_item_contents")
        for i in range(len(cards)):
            logger.debug(f"Working on review {i} for {isbn}")
            card = cards[i]
            review: data.ReviewData = data.ReviewData(
                "", LTScraper.source, "", "", "", "")
            # user
            userTag = card.select_one("a.mr_username.popup_registered")
            if userTag:
                review.username = userTag.text
            else:
                logger.critical(f"{isbn}: No username element")
                continue
            # text content
            userText = card.select_one("div.mr_review_content")
            if userText:
                review.text = userText.text
            else:
                logger.critical(f"{isbn}: No review element")
                continue
            # rating
            userRateA = card.select_one("span.rating.rating-style-306")
            if userRateA:
                userRateB = userRateA.attrs["title"]
                if userRateB:
                    review.rating = userRateB.split(" ")[0]
                else:
                    logger.critical(f"{isbn}: Invalid rating attribute")
                    continue
            else:
                logger.warn(f"{isbn}: No rating element")
                pass
            # external id
            idA = card.select_one("span.mr_note_item")
            if idA:
                idB = idA.find_next("a")
                if idB:
                    idLink = idB["href"]
                    if idLink:
                        review.external_id = idLink.split("/")[-1]
                    else:
                        logger.critical(
                            f"{isbn}: No external ID attribute")
                        continue
                else:
                    logger.critical(f"{isbn}: No external ID anchor")
                    continue
            else:
                logger.critical(f"{isbn}: No external ID span element")
                continue
            ret.append(review)
        return ret
