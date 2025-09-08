import bs4
from selenium.webdriver.common.by import By
import undetected_chromedriver as uc
import time
from . import data


class GRScraper(data.Scraper):
    source = "goodreads"
    url = "https://www.goodreads.com/book/isbn/"

    def __init__(self) -> None:
        pass

    @staticmethod
    def getReviews(isbn: str, driver: uc.Chrome) -> list[data.ReviewData]:
        ret: list[data.ReviewData] = []
        driver.get(GRScraper.url+isbn)
        # skip the GR sign-in pop-up
        time.sleep(2)
        try:
            driver.find_element(By.CLASS_NAME,
                                "Button.Button--tertiary.\
                                        Button--medium.Button--rounded")\
                .click()
        except Exception:
            pass

        soup = bs4.BeautifulSoup(
            driver.page_source, features="html.parser")
        cards = soup.select(".ReviewCard")
        # print(len(cards))
        for card in cards:
            # username
            eusername = card.select_one(".ReviewerProfile__name")
            if eusername:
                name = eusername.text
            else:
                continue

            # review text
            etextA = card.select_one(".ReviewCard__content")
            if etextA:
                etextB = etextA.select_one("span.Formatted")
                if etextB:
                    text = etextB.text
                else:
                    continue
            else:
                continue

            # rating
            eratingA = card.select_one(".ReviewCard__content")
            if eratingA:
                eratingB = eratingA.select_one(
                    "span.RatingStars.RatingStars__small")
                if eratingB:
                    # rating=PageElement(eratingB).get("aria-label")
                    eratingC = eratingB.attrs["aria-label"]
                    if eratingC:
                        rating = str(eratingC).split(" ")[1]
                    else:
                        continue
                else:
                    continue
            else:
                print("rating stage 1 not found")
                continue

            # external id
            eeidA = card.select_one(".ReviewCard__content")
            if eeidA:
                eeidB = eeidA.select_one("span.Text.Text__body3")
                if eeidB:
                    eeidC = eeidB.select_one("a")
                    if eeidC:
                        "href value of this tag is the eid"
                        try:
                            eid = str(eeidC.attrs["href"])
                        except Exception:
                            continue
                    else:
                        continue
                else:
                    continue
            else:
                continue
            # append
            ret.append(data.ReviewData("", GRScraper.source,
                       eid, name, rating, text))
        return ret
