import langchain
from langchain.chat_models import ChatOpenAI
from langchain.embeddings import SentenceTransformerEmbeddings
from langchain.agents import initialize_agent, Tool
from langchain.agents import AgentType
from webcrawler import WebCrawler
from dotenv import load_dotenv

class ChatAgent:
    def __init__(self):
        load_dotenv()
        self.llm = ChatOpenAI(model_name="gpt-3.5-turbo-16k")
        self.webcrawler = WebCrawler()
        self.tools = [
                    Tool(
                        name="websearch",
                        func=self.webcrawler.search,
                        description="useful for when you need to search for a specific topic on the web, use the query parameter to specify the what to search",
                    ),
                    # Tool(
                    #     name="generate_questions",
                    #     func=self.generate_questions,
                    #     description="useful to generate what questions to search on web using websearch tool and ask the agent",
                    # ),
                    Tool(
                        name="ask_agent",
                        func=self.ask_agent,
                        description="Useful when answers need to be rechecked. Acts like an expert in software development.",
                    ),]
        self.agent = initialize_agent(self.tools, self.llm, agent=AgentType.OPENAI_MULTI_FUNCTIONS, verbose=True)

    def ask_agent(self, question):
        answer =  self.agent.run(f"""Imagine you are a software developer,
                                 you have to help the user with their query as best as you can. You help in validating the results given to you. Do not use any tools available.
                                 Use all tools available to make your answer. Query: {question} """)
        return answer
    
    def understand_logs(self,logs):
        answer =  self.agent.run(f"""Imagine you are an expert software developer who helps in summarizing logs in high technical detail.
                                 DO NOT USE ANY TOOLS. Only give summary and no solutions. Here are the logs: \n {logs} """)
        return answer
    
    def master_agent(self,summary):
        answer =  self.agent.run(f"""Imagine you are a software developer who has to provide solutions to errors in logs of a software.
                                    Use all tools available to make your answer. Generate proper detailed question including the error and the root package in which error occured or what caused the error questions on what to ask the websearch tool.
                                    Then recheck your answer with the ask_agent tool by giving and answer and sending it to ask_agent for validation and give final solution. If websearch tool returns null reframe your question and only search the detailed log also give the links from websearch if needed.
                                    Here are the logs for which you need to find solutions: \n {summary} """)
        return answer

if __name__ == "__main__":
    agent = ChatAgent()
    summary = agent.understand_logs(""" 2023-11-28 23:45:53 Traceback (most recent call last):
        2023-11-28 23:45:53   File "/app/main.py", line 12, in <module>
        2023-11-28 23:45:53     main()
        2023-11-28 23:45:53   File "/app/main.py", line 7, in main
        2023-11-28 23:45:53     requests.get(url='http://status:8101')
        2023-11-28 23:45:53   File "/usr/local/lib/python3.9/site-packages/requests/api.py", line 73, in get
        2023-11-28 23:45:53     return request("get", url, params=params, **kwargs)
        2023-11-28 23:45:53   File "/usr/local/lib/python3.9/site-packages/requests/api.py", line 59, in request
        2023-11-28 23:45:53     return session.request(method=method, url=url, **kwargs)
        2023-11-28 23:45:53   File "/usr/local/lib/python3.9/site-packages/requests/sessions.py", line 589, in request
        2023-11-28 23:45:53     resp = self.send(prep, **send_kwargs)
        2023-11-28 23:45:53   File "/usr/local/lib/python3.9/site-packages/requests/sessions.py", line 703, in send
        2023-11-28 23:45:53     r = adapter.send(request, **kwargs)
        2023-11-28 23:45:53   File "/usr/local/lib/python3.9/site-packages/requests/adapters.py", line 519, in send
        2023-11-28 23:45:53     raise ConnectionError(e, request=request)
        2023-11-28 23:45:53 requests.exceptions.ConnectionError: HTTPConnectionPool(host='status', port=8101): Max retries exceeded with url: / (Caused by NameResolutionError("<urllib3.connection.HTTPConnection object at 0x7f4f0e703af0>: Failed to resolve 'status' ([Errno -2] Name or service not known)"))""")
    print(agent.master_agent(summary))