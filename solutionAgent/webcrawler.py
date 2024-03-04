"""Module for crawling the web."""
import os
from dotenv import load_dotenv
from tavily import TavilyClient
load_dotenv()

class WebCrawler:
    """Class for crawling the web."""
    def __init__(self):
        api_key = os.getenv('TAVILY_API_KEY')
        self.tool = TavilyClient(api_key=api_key)

    def search(self, query):
        """
        Searches for the given query using the DuckDuckGo search engine.

        Args:
            query (str): The search query.

        Returns:
            str: The first result returned by the search engine.
        """
        try:
            return(self.tool.search(query = query,include_domains = ['https://stackoverflow.com','https://github.com'],search_depth="advanced"))
            
        except Exception as e:
            print(e)
            return None
