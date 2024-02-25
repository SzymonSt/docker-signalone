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
                endpoint_url=os.environ.get("ENDPOINT_URL"),
                task="text-generation",
                model_kwargs={
                    "max_new_tokens": 512,
                    "top_k": 50,
                    "temperature": 0.4,
                    "repetition_penalty": 1.1,
                },
            )
        self.summarizer = HuggingFaceEndpoint(
                endpoint_url=os.environ.get("ENDPOINT_URL"),
                task="text-generation",
                model_kwargs={
                    "max_new_tokens": 100,
                    "top_k": 50,
                    "temperature": 0.3,
                    "repetition_penalty": 1.1,
                },
            )
        self.title_gen = HuggingFaceEndpoint(
                endpoint_url=os.environ.get("ENDPOINT_URL"),
                task="text-generation",
                model_kwargs={
                    "max_new_tokens": 60,
                    "top_k": 30,
                    "temperature": 0.5,
                    "repetition_penalty": 1.1,
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
        self.agent_executor = AgentExecutor(agent=agent, tools=self.tools, verbose=False, handle_parsing_errors=True,return_intermediate_steps=True, max_iterations=3)
     
    def understand_logs(self,logs):
        """Function to understand logs and return a summary
        Args:
            logs (str): logs from the user
        Returns: summary of the logs"""

        answer =  self.summarizer(f"""Imagine you are an expert software developer who helps in creating summary to ask websearch in detail.
                                Only give summary in form of paragraph with technical details and not solutions. Include error message in the summary. Here are the logs: \n {logs} 
                                Summary: """)
        if answer[-1] != '.':
            answer_sentences = answer.split(".")
            answer_sentences.pop()
            answer = ".".join(answer_sentences)
        return answer

    def master_agent(self,summary):
        """Function to find solutions to the logs using summary and web search
        Args:
            summary (str): summary of the logs
        Returns: solution to the logs"""

        answer =  self.agent_executor.invoke({"input":f"""Imagine you are a software developer who has to provide short summary of available solutions to issues of a software.
                                     Use websearch tool available to make your answer. You can ask multiple questions from the webagent. Use websearch agent to search about information. You can use the summary provided. Also provide the sources of your solutions.
                                     Here is the summary of issue for which you need to find solution and provide code or commands if it would help to resolve issue: \n {summary}"""})
        
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

    def generate_title(self, summary):
        """Function to generate title for the logs 
        Args:
            logs (str): logs from the user
            summary (str): summary of the logs
            solution (str): solution to the logs
            unique_id (str): unique id of the user
            userid (str): user id of the user
            container_name (str): container name of the user
        Returns: json object"""

        response = self.title_gen(f'''Give a short title for this text: {summary}.
                                  Title:''')
        # title_pattern = re.compile(r'Title:\s*(.*)', re.IGNORECASE)
        # title_match = title_pattern.search(response)
        # if title_match:
        #     title = title_match.group(1).strip()
        # else:
        #     title = "Error Log"
        return response.split("\n")[0]

    def run(self, logs):
        """Function to run the agent
        Args:
            logs (str): logs from the user
            unique_id (str): unique id of the user
            userid (str): user id of the user
            container_name (str): container name of the user
        Returns: json object"""
        summary = self.understand_logs(logs)
        urls = self.master_agent(summary)
        sol = self.llm(f'''Use this information to provide solutions to the issue summary: {summary}.
                       \n Here are the intermediate steps for you to use as in information source: {urls}
                       Provide just solutions anythign else will be punished.
                       Do not assume anything that is not there in the intermediate steps and give a proper answer.
                       Do not output any code or commands if not confirmed by the intermediate steps it must be as accurate as possible. You will be punsihed for wrong information.
                       Do not prompt user to search anything in web or ask support.
                       \n Solution:''')
        urls = self.extract_urls(logs, urls)
        title = self.generate_title(summary)
        for ui, url in enumerate(urls):
            urls[ui] = url.split("'")[0]
        final = {
            "title": title,
            "logsummary": summary,
            "predictedSolutions":sol,
            "sources": urls
                }

        return final
