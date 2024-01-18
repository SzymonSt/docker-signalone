"""Module for crawling the web."""
from duckduckgo_search import DDGS
import requests
from bs4 import BeautifulSoup

class WebCrawler:
    """Class for crawling the web."""
    def __init__(self):
        self.ddgs = DDGS()

    def search(self, query):
        """
        Searches for the given query using the DuckDuckGo search engine.

        Args:
            query (str): The search query.

        Returns:
            str: The first result returned by the search engine.
        """
        with self.ddgs as ddgs:
            results = [r["body"] for r in ddgs.text(f"{query} site:stackoverflow.com", max_results=5)]
            return(results)  

if __name__ == "__main__":
    crawler = WebCrawler()
    print(crawler.search("error in python venv"))