# this package contains data abstractions for scraper, which allow it to post
# data to the backend
from abc import abstractmethod
import undetected_chromedriver as uc


class ReviewData():
    def __init__(self, olid: str, source: str, external_id: str, username: str,
                 rating: float, text: str) -> None:
        self.olid = olid
        self.source = source
        self.external_id = external_id
        self.username = username
        self.rating = rating
        self.text = text

    def asJSON(self) -> dict[str, str | float]:
        return {
            "olid": self.olid,
            "source": self.source,
            "external_id": self.external_id,
            "username": self.username,
            "rating": self.rating,
            "text": self.text
        }


class WorkData():
    def __init__(self):
        self.olid = None
        self.title = None
        self.author = None
        self.cover = None
        self.description = None
        self.published_year = None

    def asJSON(self) -> dict[str, str]:
        return {
            "olid": self.olid,
            "title": self.title,
            "author": self.author,
            "cover": self.cover,
            "description": self.description,
            "published_year": self.published_year,
        }


class ISBNRouteData():
    def __init__(self, isbn=None, olid=None):
        self.isbn = isbn
        self.olid = olid

    def asJSON(self) -> dict[str, str]:
        return {
            "isbn": self.isbn,
            "olid": self.olid,
        }


class Scraper():
    source: str = "null"
    url = ""
    interval: tuple[int, int] = (0, 0)

    @staticmethod
    @abstractmethod
    def getReviews(isbn: str, driver: uc.Chrome) -> list[ReviewData]:
        return []
