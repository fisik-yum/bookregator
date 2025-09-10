from olclient import OpenLibrary
import isbnlib as islib
from . import data

global ol
ol = OpenLibrary()


def ISBNtoOLIDW(isbn_val: str) -> (str, str):
    if not islib.is_isbn10(isbn_val) and not islib.is_isbn13(isbn_val):
        raise Exception("Invalid ISBN")
    isbn_val = islib.to_isbn13(isbn_val)
    return (isbn_val, ol.Edition.get(isbn=isbn_val).work_olid)


def GenerateWorkData(olidw: str) -> data.WorkData:
    work = data.WorkData()

    temp = ol.Work.get(olidw)

    work.olid = olidw
    work.description = temp.description
    work.title = temp.title
    if hasattr(temp,"covers"):
        work.cover = temp.covers[0] if len(temp.covers) > 0 else ""
    # TODO: fix author
    # work.author = temp.authors
    return work
