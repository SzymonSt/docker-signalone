import json
import langchain
from langchain.chat_models import ChatOpenAI
from langchain.embeddings import SentenceTransformerEmbeddings
from langchain.agents import initialize_agent, Tool
from langchain.agents import AgentType
from webcrawler import WebCrawler
from dotenv import load_dotenv
from openai import OpenAI
import time

class ChatAgent:
    def __init__(self):
        load_dotenv()
        self.llm = ChatOpenAI(model_name="gpt-3.5-turbo-16k", temperature=0.5)
        self.client = OpenAI()

        self.webcrawler = WebCrawler()
        self.tools = [
                    Tool(
                        name="websearch",
                        func=self.webcrawler.search,
                        description="useful for when you need to search for a specific topic on the web, use the query parameter to specify the what to search",
                    ),
                    ]
        self.agent = initialize_agent(self.tools, self.llm, agent=AgentType.OPENAI_MULTI_FUNCTIONS, verbose=True)

    def understand_logs(self,logs):
        answer =  self.agent.run(f"""Imagine you are an expert software developer who helps in summarizing logs in high technical detail in points.
                                 DO NOT USE ANY TOOLS. Only give summary and no solutions. Here are the logs: \n {logs} """)
        return answer
    
    def master_agent(self,summary):
        answer =  self.agent.run(f"""Imagine you are a software developer who has to provide solutions to errors in logs of a software.
                                     Use all tools available to make your answer. You can use the summary provided. Also provide the sources of your solutions.
                                     Here are the logs for which you need to find solutions: \n {summary}""")
        return answer
    
    def generate_json(self, logs, summary, solution, unique_id="8uueerwer8248289", userid="iwuiji3h3928r7hyuwr", container_name="testcontainer"):
        response = self.client.chat.completions.create(
        model="gpt-3.5-turbo-1106",
        response_format={ "type": "json_object" },
        messages=[
            {"role": "system", "content": "You are a helpful assistant designed to output JSON."},
            {"role": "user", "content": f""" Use {solution} to fill in the following JSON object values for "sources"
             array should include all links and sources. Give an appropriate title to the logs and summary in "title"
                            For your answer create a JSON object with the following schema:
                            {{ "id": "unique_id", 
                             "userid": "userid",
                              "containerName": "container_name",
                             "logs":{logs},
                              "title": "Give a title",
                               "timestamp": {time.time()},
                               "logsummary": {summary},
                               "predicted_solution":{solution},
                                "sources": ["source1","source2","etc"]
                                }} only return json object and nothing else.
                            json_object:"""}
        ]
        )
        return(response.choices[0].message.content)
    
    def run(self, logs, unique_id="test_id", userid="test_user1", container_name="testcontainer"):
        summary = self.understand_logs(logs)
        solution = self.master_agent(summary)
        return self.generate_json(logs, summary, solution)
    
if __name__ == "__main__":
    agent = ChatAgent()
    logs = """ 2023-11-28 23:45:53 Traceback (most recent call last):
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
        2023-11-28 23:45:53 requests.exceptions.ConnectionError: HTTPConnectionPool(host='status', port=8101): Max retries exceeded with url: / (Caused by NameResolutionError("<urllib3.connection.HTTPConnection object at 0x7f4f0e703af0>: Failed to resolve 'status' ([Errno -2] Name or service not known)"))"""
    
    print(agent.run(logs))
