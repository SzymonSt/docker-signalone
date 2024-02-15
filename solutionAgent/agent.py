"""Module for the chat agent."""
import json
import os
import re
from langchain import hub
from langchain_community.llms import HuggingFaceEndpoint
from langchain.agents import initialize_agent, Tool, create_react_agent, AgentExecutor
from langchain.agents import AgentType
from webcrawler import WebCrawler
from dotenv import load_dotenv
from openai import OpenAI
from datetime import datetime

class ChatAgent:
    """Class for the chat agent."""
    def __init__(self):
        load_dotenv()
        self.llm = HuggingFaceEndpoint(
                endpoint_url=os.getenv("ENDPOINT_URL"),
                task="text-generation",
                model_kwargs={
                    "max_new_tokens": 512,
                    "top_k": 50,
                    "temperature": 0.5,
                    "repetition_penalty": 1.03,
                },
            )
        self.webcrawler = WebCrawler()
        self.tools = [
                    Tool(
                        name="websearch",
                        func=self.webcrawler.search,
                        description="useful for when you need to search for a specific topic on the web, use the query parameter to specify the what to search",
                    )
                    ]
        
        prompt = hub.pull("hwchase17/react")
        agent = create_react_agent(self.llm, self.tools, prompt)
        self.agent_executor = AgentExecutor(agent=agent, tools=self.tools, verbose=True, handle_parsing_errors=True,return_intermediate_steps=True, max_iterations=2) 
     
    def understand_logs(self,logs):
        """Function to understand logs and return a summary
        Args:
            logs (str): logs from the user
        Returns: summary of the logs"""

        answer =  self.llm(f"""Imagine you are an expert software developer who helps in creating summary to ask websearch in detail.
                                Only give summary in form of paragraph and not solutions. Here are the logs: \n {logs} 
                                Summary: """)
        return answer

    def master_agent(self,summary):
        """Function to find solutions to the logs using summary and web search
        Args:
            summary (str): summary of the logs
        Returns: solution to the logs"""

        answer =  self.agent_executor.invoke({"input":f"""Imagine you are a software developer who has to provide short summary of available solutions to errors in logs of a software.
                                     Use all tools available to make your answer. You can ask multiple questions from the webagent. Use websearch agent to search about information. You can use the summary provided. Also provide the sources of your solutions.
                                     Here are the logs for which you need to find solution and provide code if necessary: \n {summary}"""})
        
        return answer['intermediate_steps']
    
    def extract_urls(self,logs, urls):
        """Function to extract urls from the intermediate steps
        Args:
            urls (list): list of urls
        Returns: list of urls"""

        prompt = f'''You are a helpful assistant designed to give formatted output. Format to be followed [url1,url2,url3, .....].
        Extract only the urls in proper format from the intermediate steps and return them in a list.
        Only extract the relevant urls which solve the error in these logs {logs} \n Here are the intermediate steps: {urls}\n
        Response: List(urls) ='''
        urls = self.llm(prompt)
        urls = re.findall(r'https?://\S+', urls)
        return urls

    def generate_severity(self, logs):
        """Function to generate json object 
        Args:
            logs (str): logs from the user
            summary (str): summary of the logs
            solution (str): solution to the logs
            unique_id (str): unique id of the user
            userid (str): user id of the user
            container_name (str): container name of the user
        Returns: json object"""

        response = self.llm(f'''Give a relevant title and severity for these logs: {logs}.
                            Only give the title and severity and nothing else in your output.
                            Give output in a structured format like:
                            \n Title: <title for the error log>\n
                                Severity: <severity of the error log>\n
                            Do not include any other parameters in the output.''')

        title_pattern = re.compile(r'Title:\s*(.*)', re.IGNORECASE)
        severity_pattern = re.compile(r'Severity:\s*(\w+)', re.IGNORECASE)

        title_match = title_pattern.search(response)
        severity_match = severity_pattern.search(response)

        if title_match:
            title = title_match.group(1).strip()
        else:
            title = "Error Log"

        if severity_match:
            severity = severity_match.group(1).strip()
        else:
            severity = "High"
        return title, severity

    def run(self, logs, unique_id="test_id", userid="test_user1", container_name="testcontainer"):
        """Function to run the agent
        Args:
            logs (str): logs from the user
            unique_id (str): unique id of the user
            userid (str): user id of the user
            container_name (str): container name of the user
        Returns: json object"""
        summary = self.understand_logs(logs)
        urls = self.master_agent(summary)
        sol = self.llm(f'''Use this information to provide the solution to the logs: {summary}.
                       \n Here are the intermediate steps for you to use as in information source: {urls}
                       Do not assume anything that is not there in the intermediate steps and give a proper answer.
                       \n Solution:''')
        urls = self.extract_urls(logs, urls)
        title , severity = self.generate_severity(logs)        
        final = {
            "id": {unique_id},
            "userid": userid,
            "containerName": container_name,
            "logs":logs,
            "severity": severity,
            "title": title,
            "timestamp": str(datetime.now()),
            "logsummary": summary,
            "predicted_solution":sol,
            "sources": urls
                }

        return final
