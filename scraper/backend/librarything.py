import bs4
from selenium.webdriver.common.by import By
import undetected_chromedriver as uc
import time
from . import data


class LTScraper(data.Scraper):
    url = "https://www.librarything.com/isbn/"

    @staticmethod
    def ISBNtoLTCode(isbn: str, driver: uc.Chrome) -> str:
        searchurl: str = "https://www.librarything.com/search.php?search=" + \
            isbn+"&searchtype=newwork_titles&sortchoice=0"

        driver.get(searchurl)
        soup = bs4.BeautifulSoup(driver.page_source, features="html.parser")
        table = soup.find(
            "table", attrs={"class": "msg", "border": "0", "cellspacing": "0"})
        if table:
            return table.a["data-workid"].__str__()

    @staticmethod
    def getReviews(isbn: str, driver: uc.Chrome) -> list[data.ReviewData]:
        ret: list[data.ReviewData] = []
        driver.get(LTScraper.url+isbn)
        time.sleep(10)

        return []
