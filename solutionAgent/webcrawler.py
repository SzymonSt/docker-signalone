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
            return(self.tool.search(query = query,include_domains = ['https://github.com','https://stackoverflow.com'],search_depth="advanced"))
            
        except Exception as e:
            print(e)
            return None

if __name__ == '__main__':
    crawl = WebCrawler()
    print(crawl.search("""AttributeError: module 'numpy' has no attribute 'object'.
`np.object` was a deprecated alias for the builtin `object`. To avoid this error in existing code, use `object` by itself. Doing this will not modify any behavior and is safe. 
The aliases was originally deprecated in NumPy 1.20; for more details and guidance see the original release note at:
    https://numpy.org/devdocs/release/1.20.0-notes.html#deprecations"""))