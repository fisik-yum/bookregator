from olclient import OpenLibrary
import isbnlib as islib
import requests as _requests
from . import data

global ol
ol = OpenLibrary()

_OL_API = "https://openlibrary.org"

_NON_BOOK_FORMATS = {
    "cd", "audio cd", "mp3 cd", "dvd", "vinyl", "lp",
    "audio cassette", "cassette", "audible audio", "audio download",
}


def ISBNtoOLIDW(isbn_val: str) -> (str, str):
    isbn_val = islib.canonical(isbn_val)
    if not islib.is_isbn10(isbn_val) and not islib.is_isbn13(isbn_val):
        raise Exception("Invalid ISBN")
    isbn_val = islib.to_isbn13(isbn_val)
    try:
        edition = ol.Edition.get(isbn=isbn_val)
        fmt = (getattr(edition, "physical_format", "") or "").lower().strip()
        if fmt in _NON_BOOK_FORMATS:
            raise Exception(f"Non-book format: {fmt}")
        identifiers = getattr(edition, "identifiers", {}) or {}
        if "music_brainz" in identifiers or "musicbrainz" in identifiers:
            raise Exception("MusicBrainz edition, skipping")
        return isbn_val, edition.work_olid
    except Exception:
        raise Exception("Invalid Attr")


def GenerateWorkData(olidw: str) -> data.WorkData:
    resp = _requests.get(f"{_OL_API}/works/{olidw}.json", timeout=10)
    resp.raise_for_status()
    j = resp.json()

    work_type = j.get("type", {}).get("key", "")
    if work_type and work_type != "/type/work":
        raise Exception(f"Unexpected OL type {work_type} for {olidw}")

    work = data.WorkData()
    work.olid = olidw
    work.title = j.get("title", "")
    desc = j.get("description", "")
    work.description = desc.get("value", desc) if isinstance(desc, dict) else desc
    covers = j.get("covers", [])
    work.cover = str(covers[0]) if covers else ""
    return work
