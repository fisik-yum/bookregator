# uv run scrapy runspider spider.py
import asyncio
import json
import logging
import threading

import bs4
import requests
import scrapy
from scrapy.crawler import CrawlerProcess
from scrapy_playwright.page import PageMethod
from playwright_stealth import Stealth
import undetected_chromedriver as uc

from backend import data, olshim
from backend.goodreads import GRScraper
from backend.librarything import LTScraper

logger = logging.getLogger(__name__)

AUTH = ("admin", "opensesame")
API = "http://127.0.0.1:1024/api"

GR_CONTEXT = {
    "user_agent": (
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
        "AppleWebKit/537.36 (KHTML, like Gecko) "
        "Chrome/139.0.0.0 Safari/537.36"
    ),
    "viewport": {"width": 1280, "height": 900},
    "locale": "en-US",
}


async def _apply_stealth(page, request):
    await Stealth().apply_stealth_async(page)


def _playwright_meta(**extra):
    return {
        "playwright": True,
        "playwright_context": "default",
        "playwright_page_init_callback": _apply_stealth,
        "playwright_page_methods": [PageMethod("wait_for_timeout", 4000)],
        **extra,
    }


class BookSpider(scrapy.Spider):
    name = "bookspider"

    custom_settings = {
        "DOWNLOAD_HANDLERS": {
            "https": "scrapy_playwright.handler.ScrapyPlaywrightDownloadHandler",
            "http": "scrapy_playwright.handler.ScrapyPlaywrightDownloadHandler",
        },
        "PLAYWRIGHT_LAUNCH_OPTIONS": {
            "headless": True,
            "args": [
                "--no-sandbox",
                "--disable-blink-features=AutomationControlled",
            ],
        },
        "PLAYWRIGHT_CONTEXTS": {"default": GR_CONTEXT},
        "PLAYWRIGHT_PROCESS_REQUEST_HEADERS": None,
        "TWISTED_REACTOR": "twisted.internet.asyncioreactor.AsyncioSelectorReactor",
        "CONCURRENT_REQUESTS": 4,
        "DOWNLOAD_DELAY": 2,
        "AUTOTHROTTLE_ENABLED": True,
        "LOG_LEVEL": "INFO",
    }

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._seen_isbns: set[str] = set()
        self._seen_olids: set[str] = set()
        # UC driver for LT — Cloudflare blocks headless Playwright on LT search pages.
        # Threading lock serializes access since Selenium is not thread-safe.
        self._lt_driver = uc.Chrome(use_subprocess=False, version_main=148)
        self._lt_lock = threading.Lock()

    def closed(self, reason):
        self._lt_driver.quit()

    # ── seed loading ────────────────────────────────────────────────────────

    def start_requests(self):
        yield from self._seed_isbn("978-0593135204")

    def _seed_isbn(self, raw_isbn: str):
        """Resolve ISBN → OLID, post work/route metadata, yield a GR Playwright request."""
        if raw_isbn in self._seen_isbns:
            return
        try:
            isbn, olid = olshim.ISBNtoOLIDW(raw_isbn)
        except Exception:
            return
        if isbn in self._seen_isbns:
            return
        self._seen_isbns.add(isbn)

        if olid not in self._seen_olids:
            try:
                work = olshim.GenerateWorkData(olid)
                route = data.ISBNRouteData(isbn, olid)
                requests.post(f"{API}/insert/work", json=work.asJSON(), auth=AUTH)
                requests.post(f"{API}/insert/route", json=route.asJSON(), auth=AUTH)
                self._seen_olids.add(olid)
            except Exception as e:
                logger.warning(f"Skipping {olid} ({isbn}): {e}")
                return

        yield scrapy.Request(
            f"https://www.goodreads.com/book/isbn/{isbn}",
            callback=self.parse_gr,
            errback=self.handle_error,
            meta=_playwright_meta(isbn=isbn, olid=olid),
        )

    # ── GR scraping ─────────────────────────────────────────────────────────

    async def parse_gr(self, response):
        isbn = response.meta["isbn"]
        olid = response.meta["olid"]

        soup = bs4.BeautifulSoup(response.text, "html.parser")
        review_cards = soup.select(".ReviewCard")
        carousel_links = response.css(".CarouselGroup a[href*='/book/show']").getall()
        logger.info(
            f"GR {isbn}: status={response.status} len={len(response.text)} "
            f"reviews={len(review_cards)} carousel={len(carousel_links)}"
        )

        gr_reviews = GRScraper._parse_reviews(soup, isbn)
        for r in gr_reviews:
            r.olid = olid

        loop = asyncio.get_event_loop()

        if gr_reviews:
            await loop.run_in_executor(
                None,
                lambda: requests.post(
                    f"{API}/insert/reviewmultiple",
                    json=[r.asJSON() for r in gr_reviews],
                    auth=AUTH,
                ),
            )

        # LT runs in a thread — serialized by lock so the UC driver isn't shared.
        await self._scrape_lt(isbn, olid, loop)

        # Discover new ISBNs from "Readers also enjoyed" carousel.
        for href in response.css(
            ".CarouselGroup a[href*='/book/show']::attr(href)"
        ).getall():
            yield scrapy.Request(
                href.split("?")[0],
                callback=self.parse_gr_discover,
                errback=self.handle_error,
                meta=_playwright_meta(),
            )

    async def parse_gr_discover(self, response):
        """Follow a /book/show/ link, pull the ISBN out of JSON-LD, re-seed."""
        for script_text in response.css(
            'script[type="application/ld+json"]::text'
        ).getall():
            try:
                schema = json.loads(script_text)
                if schema.get("@type") == "Book" and "isbn" in schema:
                    for req in self._seed_isbn(schema["isbn"]):
                        yield req
                    break
            except (json.JSONDecodeError, KeyError):
                continue

    async def _scrape_lt(self, isbn: str, olid: str, loop: asyncio.AbstractEventLoop):
        def _run():
            with self._lt_lock:
                return LTScraper.getReviews(isbn, self._lt_driver)

        try:
            lt_reviews = await loop.run_in_executor(None, _run)
            for r in lt_reviews:
                r.olid = olid
            if lt_reviews:
                await loop.run_in_executor(
                    None,
                    lambda: requests.post(
                        f"{API}/insert/reviewmultiple",
                        json=[r.asJSON() for r in lt_reviews],
                        auth=AUTH,
                    ),
                )
        except Exception as e:
            logger.error(f"LT scrape failed for {isbn}: {e}")

    def handle_error(self, failure):
        logger.error(f"Request failed: {failure.request.url} — {failure.value}")


if __name__ == "__main__":
    process = CrawlerProcess()
    process.crawl(BookSpider)
    process.start()
