import time
from fastapi import FastAPI
from pydantic import BaseModel
from agent import ChatAgent

class LogData(BaseModel):
    '''Class for the log data'''
    logs: str
    unique_id: str = "test_id"
    userid: str = "test_user1"
    container_name: str = "testcontainer"

app = FastAPI()

@app.post("/run_chat_agent")
async def run_chat_agent(data: LogData):
    '''Function to run the chat agent'''
    chat_agent = ChatAgent()
    retries = 0
    while True:
        try:
            print(f"Number of retries {retries}")
            retries =retries + 1
            result = chat_agent.run(data.logs, data.unique_id, data.userid, data.container_name)
            return result
        except Exception as e:
            if retries > 4:
                return {"error": "Unable to process the logs"}
            time.sleep(5)
        