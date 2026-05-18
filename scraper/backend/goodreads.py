import bs4
from selenium.webdriver.common.by import By
import undetected_chromedriver as uc
import time
import coloredlogs
import logging
from . import data

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
coloredlogs.install()


class GRScraper(data.Scraper):
    source = "goodreads"
    url = "https://www.goodreads.com/book/isbn/"

    def __init__(self) -> None:
        pass

    @staticmethod
    def _parse_reviews(soup: bs4.BeautifulSoup, isbn: str) -> list[data.ReviewData]:
        """Parse reviews from already-rendered BeautifulSoup. Called by getReviews and the Scrapy spider."""
        ret: list[data.ReviewData] = []
        logger = logging.getLogger(__name__)
        cards = soup.select(".ReviewCard")
        for i, card in enumerate(cards):
            logger.debug(f"Working on review {i} for {isbn}")

            eusername = card.select_one(".ReviewerProfile__name")
            if eusername:
                name = eusername.text
            else:
                logger.critical(f"{isbn}: No username element")
                continue

            etextA = card.select_one(".ReviewCard__content")
            if etextA:
                etextB = etextA.select_one("span.Formatted")
                if etextB:
                    for br in etextB.find_all("br"):
                        br.replace_with("\n")
                    text = etextB.text
                else:
                    logger.critical(f"{isbn}: No review span")
                    continue
            else:
                logger.critical(f"{isbn}: No review element")
                continue

            eratingA = card.select_one(".ReviewCard__content")
            if eratingA:
                eratingB = eratingA.select_one("span.RatingStars.RatingStars__small")
                if eratingB:
                    eratingC = eratingB.attrs.get("aria-label", "")
                    if eratingC:
                        rating = float(eratingC.split(" ")[1])
                    else:
                        logger.critical(f"{isbn}: Invalid rating attribute")
                        continue
                else:
                    logger.warning(f"{isbn}: No rating span")
                    continue
            else:
                logger.warning(f"{isbn}: No rating element")
                continue

            eeidA = card.select_one(".ReviewCard__content")
            if eeidA:
                eeidB = eeidA.select_one("span.Text.Text__body3")
                if eeidB:
                    eeidC = eeidB.select_one("a")
                    if eeidC:
                        try:
                            eid = str(eeidC.attrs["href"])
                        except Exception:
                            continue
                    else:
                        logger.critical(f"{isbn}: No external ID attribute")
                        continue
                else:
                    logger.critical(f"{isbn}: No external ID anchor")
                    continue
            else:
                logger.critical(f"{isbn}: No external ID span element")
                continue

            ret.append(data.ReviewData("", GRScraper.source, eid, name, rating, text))
        return ret

    @staticmethod
    def getReviews(isbn: str, driver: uc.Chrome) -> list[data.ReviewData]:
        logger = logging.getLogger(__name__)
        logger.debug(f"Fetching details for {isbn}")
        driver.get(GRScraper.url + isbn)
        time.sleep(2)
        try:
            driver.find_element(
                By.CLASS_NAME,
                "Button.Button--tertiary.Button--medium.Button--rounded"
            ).click()
        except Exception:
            pass
        soup = bs4.BeautifulSoup(driver.page_source, features="html.parser")
        return GRScraper._parse_reviews(soup, isbn)
