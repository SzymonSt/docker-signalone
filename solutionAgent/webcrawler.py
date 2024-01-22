"""Module for crawling the web."""
from duckduckgo_search import DDGS
import requests
from bs4 import BeautifulSoup
from langchain_community.utilities import StackExchangeAPIWrapper
from langchain_community.tools.tavily_search import TavilySearchResults


class WebCrawler:
    """Class for crawling the web."""
    def __init__(self):
        self.ddgs = DDGS()
        self.stackexchange = StackExchangeAPIWrapper()
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
            # with self.ddgs as ddgs:
            #     results = [r["body"] for r in ddgs.text(f"{query} site:stackoverflow.com", max_results=5)]
            #     return(results)
            return(self.tool.invoke({"query": query},search_depth="advanced"))
            
        except Exception as e:
            print(e)
            return None
          

if __name__ == "__main__":
    wc = WebCrawler()
    print(wc.search(r'''InvalidRequestError Traceback (most recent call last) Cell In[40], line 2 1 text = "What would be a good company name for a company that makes colorful socks?" ----> 2 print(llm(text) '''))
#     print(wc.search(r'''InvalidRequestError Traceback (most recent call last) Cell In[40], line 2 1 text = "What would be a good company name for a company that makes colorful socks?" ----> 2 print(llm(text) \

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\langchain\llms\base.py:291, in BaseLLM.call(self, prompt, stop, callbacks) 286 def call( 287 self, prompt: str, stop: Optional[List[str]] = None, callbacks: Callbacks = None 288 ) -> str: 289 """Check Cache and run the LLM on the given prompt and input.""" 290 return ( --> 291 self.generate([prompt], stop=stop, callbacks=callbacks) 292 .generations[0][0] 293 .text 294 ) \

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\langchain\llms\base.py:191, in BaseLLM.generate(self, prompts, stop, callbacks) 189 except (KeyboardInterrupt, Exception) as e: 190 run_manager.on_llm_error(e) --> 191 raise e 192 run_manager.on_llm_end(output) 193 return output

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\langchain\llms\base.py:185, in BaseLLM.generate(self, prompts, stop, callbacks) 180 run_manager = callback_manager.on_llm_start( 181 {"name": self.class.name}, prompts, invocation_params=params 182 ) 183 try: 184 output = ( --> 185 self._generate(prompts, stop=stop, run_manager=run_manager) 186 if new_arg_supported 187 else self._generate(prompts, stop=stop) 188 ) 189 except (KeyboardInterrupt, Exception) as e: 190 run_manager.on_llm_error(e)

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\langchain\llms\openai.py:315, in BaseOpenAI._generate(self, prompts, stop, run_manager) 313 choices.extend(response["choices"]) 314 else: --> 315 response = completion_with_retry(self, prompt=_prompts, **params) 316 choices.extend(response["choices"]) 317 if not self.streaming: 318 # Can't update token usage if streaming

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\langchain\llms\openai.py:106, in completion_with_retry(llm, **kwargs) 102 @retry_decorator 103 def _completion_with_retry(**kwargs: Any) -> Any: 104 return llm.client.create(**kwargs) --> 106 return _completion_with_retry(**kwargs)

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\tenacity_init_.py:289, in BaseRetrying.wraps..wrapped_f(*args, **kw) 287 @functools.wraps(f) 288 def wrapped_f(*args: t.Any, **kw: t.Any) -> t.Any: --> 289 return self(f, *args, **kw)

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\tenacity_init_.py:379, in Retrying.call(self, fn, *args, **kwargs) 377 retry_state = RetryCallState(retry_object=self, fn=fn, args=args, kwargs=kwargs) 378 while True: --> 379 do = self.iter(retry_state=retry_state) 380 if isinstance(do, DoAttempt): 381 try:

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\tenacity_init_.py:314, in BaseRetrying.iter(self, retry_state) 312 is_explicit_retry = fut.failed and isinstance(fut.exception(), TryAgain) 313 if not (is_explicit_retry or self.retry(retry_state)): --> 314 return fut.result() 316 if self.after is not None: 317 self.after(retry_state)

# File ~\AppData\Local\Programs\Python\Python39\lib\concurrent\futures_base.py:438, in Future.result(self, timeout) 436 raise CancelledError() 437 elif self._state == FINISHED: --> 438 return self.__get_result() 440 self._condition.wait(timeout) 442 if self._state in [CANCELLED, CANCELLED_AND_NOTIFIED]:

# File ~\AppData\Local\Programs\Python\Python39\lib\concurrent\futures_base.py:390, in Future.__get_result(self) 388 if self._exception: 389 try: --> 390 raise self._exception 391 finally: 392 # Break a reference cycle with the exception in self._exception 393 self = None

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\tenacity_init_.py:382, in Retrying.call(self, fn, *args, **kwargs) 380 if isinstance(do, DoAttempt): 381 try: --> 382 result = fn(*args, **kwargs) 383 except BaseException: # noqa: B902 384 retry_state.set_exception(sys.exc_info()) # type: ignore[arg-type]

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\langchain\llms\openai.py:104, in completion_with_retry.._completion_with_retry(**kwargs) 102 @retry_decorator 103 def _completion_with_retry(**kwargs: Any) -> Any: --> 104 return llm.client.create(**kwargs)

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\openai\api_resources\completion.py:25, in Completion.create(cls, *args, **kwargs) 23 while True: 24 try: ---> 25 return super().create(*args, **kwargs) 26 except TryAgain as e: 27 if timeout is not None and time.time() > start + timeout:

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\openai\api_resources\abstract\engine_api_resource.py:149, in EngineAPIResource.create(cls, api_key, api_base, api_type, request_id, api_version, organization, **params) 127 @classmethod 128 def create( 129 cls, (...) 136 **params, 137 ): 138 ( 139 deployment_id, 140 engine, 141 timeout, 142 stream, 143 headers, 144 request_timeout, 145 typed_api_type, 146 requestor, 147 url, 148 params, --> 149 ) = cls.__prepare_create_request( 150 api_key, api_base, api_type, api_version, organization, **params 151 ) 153 response, _, api_key = requestor.request( 154 "post", 155 url, (...) 160 request_timeout=request_timeout, 161 ) 163 if stream: 164 # must be an iterator

# File ~\AppData\Local\Programs\Python\Python39\lib\site-packages\openai\api_resources\abstract\engine_api_resource.py:83, in EngineAPIResource.__prepare_create_request(cls, api_key, api_base, api_type, api_version, organization, **params) 81 if typed_api_type in (util.ApiType.AZURE, util.ApiType.AZURE_AD): 82 if deployment_id is None and engine is None: ---> 83 raise error.InvalidRequestError( 84 "Must provide an 'engine' or 'deployment_id' parameter to create a %s" 85 % cls, 86 "engine", 87 ) 88 else: 89 if model is None and engine is None:

# InvalidRequestError: Must provide an 'engine' or 'deployment_id' parameter to create a <class 'openai.api_resources.completion.Completion'> '''))
    

    

    