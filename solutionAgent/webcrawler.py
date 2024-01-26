"""Module for crawling the web."""
from langchain_community.tools.tavily_search import TavilySearchResults


class WebCrawler:
    """Class for crawling the web."""
    def __init__(self):
        self.tool = TavilySearchResults()

    def search(self, query):
        """
        Searches for the given query using the DuckDuckGo search engine.

        Args:
            query (str): The search query.

        Returns:
            str: The first result returned by the search engine.
        """
        try:
            return(self.tool.invoke({"query": query},search_depth="advanced"))
            
        except Exception as e:
            print(e)
            return None
        