from abc import abstractmethod
import undetected_chromedriver as uc


class ReviewData():
    def __init__(self, olid: str, source: str, external_id: str, username: str,
                 rating: str, text: str) -> None:
        self.olid = olid
        self.source = source
        self.external_id = external_id
        self.username = username
        self.rating = rating
        self.text = text

    def asJSON(self) -> dict[str, str]:
        return {
            "olid": self.olid,
            "source": self.source,
            "external_id": self.external_id,
            "username": self.username,
            "rating": self.rating,
            "text": self.text
        }


class Scraper():
    source: str = "null"
    url = ""
    interval: tuple[int, int] = (0, 0)

    @staticmethod
    @abstractmethod
    def getReviews(isbn: str, driver: uc.Chrome) -> list[ReviewData]:
        return []
