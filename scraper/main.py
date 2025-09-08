from abc import abstractmethod
from selenium.webdriver.common.by import By
import undetected_chromedriver as uc
import bs4
from undetected_chromedriver.cdp import PageElement, requests
import time
driver = uc.Chrome(use_subprocess=False)

class ReviewData():
    def __init__(self, olid: str, source: str, external_id: str, username: str,
                 rating: str, text: str) -> None:
        self.olid=olid
        self.source=source
        self.external_id=external_id
        self.username=username
        self.rating=rating
        self.text=text
    def asJSON(self)->dict[str,str]:
        return {
                    "olid":self.olid,
                    "source": self.source,
                    "external_id": self.external_id,
                    "username": self.username,
                    "rating": self.rating,
                    "text": self.text
                }
class Scraper():
    source: str = "null"
    url=""
    interval: tuple[int, int] = (0,0)
    @staticmethod
    @abstractmethod
    def getReviews(isbn: str)->list[ReviewData]:
        return []


class GRScraper(Scraper):
    source="goodreads"
    url = "https://www.goodreads.com/book/isbn/"
    def __init__(self)-> None:
        pass
    @staticmethod
    def getReviews(isbn: str) -> list[ReviewData]:
        ret: list[ReviewData]=[]
        driver.get(GRScraper.url+isbn)
        # skip the GR sign-in pop-up
        time.sleep(2)
        try:
            driver.find_element(By.CLASS_NAME,
                "Button.Button--tertiary.Button--medium.Button--rounded").click()
        except:
            pass

        soup=bs4.BeautifulSoup(driver.page_source, features="html.parser")
        cards=soup.select(".ReviewCard")
        #print(len(cards))
        for card in cards:
            # username
            eusername=card.select_one(".ReviewerProfile__name")
            if eusername:
                name=eusername.text
            else:
                continue

            # review text
            etextA=card.select_one(".ReviewCard__content")
            if etextA:
                etextB=etextA.select_one("span.Formatted")
                if etextB:
                    text=etextB.text
                else:
                    continue
            else: 
                continue

            # rating
            eratingA=card.select_one(".ReviewCard__content")
            if eratingA:
                eratingB=eratingA.select_one("span.RatingStars.RatingStars__small")
                if eratingB:
                    #rating=PageElement(eratingB).get("aria-label")
                    eratingC=eratingB.attrs["aria-label"]
                    if eratingC:
                        rating=str(eratingC).split(" ")[1]
                    else:
                        continue
                else:
                    continue
            else:
                print("rating stage 1 not found")
                continue

            # external id
            eeidA=card.select_one(".ReviewCard__content")
            if eeidA:
                eeidB=eeidA.select_one("span.Text.Text__body3")
                if eeidB:
                    eeidC=eeidB.select_one("a")
                    if eeidC:
                        "href value of this tag is the eid"
                        try:
                            eid=str(eeidC.attrs["href"])
                        except:
                            continue
                    else:
                        continue
                else:
                    continue
            else:
                continue
            # append
            ret.append(ReviewData("",GRScraper.source,eid,name,rating,text))
        return ret

class LTScraper(Scraper):
    url="https://www.librarything.com/isbn/"

    @staticmethod
    def getReviews(isbn: str)->list[ReviewData]:
        ret: list[ReviewData]=[]
        driver.get(LTScraper.url+isbn)
        time.sleep(10)

        return []
def main():
    print("Getting reviews")
    revs=GRScraper().getReviews("9780375822070")
    print(f"# of reviews: {len(revs)}")
    print("posting reviews")
    for review in revs:
        review.olid="OL45804"

        requests.post("http://127.0.0.1:1024/api/insert/reviewsingle", json=review.asJSON())
main()
