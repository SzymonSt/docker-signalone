"""Module for the chat agent."""
import json
from langchain.chat_models import ChatOpenAI
from langchain.agents import initialize_agent, Tool
from langchain.agents import AgentType
from webcrawler import WebCrawler
from dotenv import load_dotenv
from openai import OpenAI
from datetime import datetime

class ChatAgent:
    """Class for the chat agent."""
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
        """Function to understand logs and return a summary
        Args:
            logs (str): logs from the user
        Returns: summary of the logs"""

        answer =  self.agent.run(f"""Imagine you are an expert software developer who helps in summarizing logs in high technical detail in points.
                                 DO NOT USE ANY TOOLS. Only give summary and no solutions. Here are the logs: \n {logs} """)
        return answer

    def master_agent(self,summary):
        """Function to find solutions to the logs using summary and web search
        Args:
            summary (str): summary of the logs
        Returns: solution to the logs"""

        answer =  self.agent.run(f"""Imagine you are a software developer who has to provide solutions to errors in logs of a software.
                                     Use all tools available to make your answer. You can use the summary provided. Also provide the sources of your solutions.
                                     Here are the logs for which you need to find solutions: \n {summary}""")
        return answer
    
    def generate_json(self, logs, summary, solution, unique_id, userid, container_name):
        """Function to generate json object 
        Args:
            logs (str): logs from the user
            summary (str): summary of the logs
            solution (str): solution to the logs
            unique_id (str): unique id of the user
            userid (str): user id of the user
            container_name (str): container name of the user
        Returns: json object"""

        response = self.client.chat.completions.create(
        model="gpt-3.5-turbo-1106",
        response_format={ "type": "json_object" },
        messages=[
            {"role": "system", "content": "You are a helpful assistant designed to output JSON."},
            {"role": "user", "content": f""" Use {solution} to fill in the following JSON object values for "sources"
             array should include all links and sources. Give an appropriate title to the logs and summary in "title"
                            For your answer create a JSON object with the following schema:
                            {{ "id": {unique_id}, 
                             "userid": {userid},
                              "containerName": {container_name},
                             "logs":{logs},
                              "title": "Give a title",
                               "timestamp": { str(datetime.now())},
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
        return json.loads(self.generate_json(logs, summary, solution, unique_id, userid, container_name))

