import requests

url = "http://localhost:8000/run_chat_agent"
data = {
    "logs": "ValueError: Error raised by inference API: Model HuggingFaceH4/zephyr-7b-beta is currently loading      ",
    "unique_id": "32234123",
    "userid": "sdfwff23f32",
    "container_name": "232332ewe"
}
response = None
while response is None:
    response = requests.post(url, json=data)

print(response.json())