import requests
from qdrant_client import QdrantClient
import json
import re
import datetime

version = "0.0.1"
PROMPT_LOG_PATH = "../prompt.log"
INFERENCE_BASE_URL = "https://api-inference.huggingface.co/models/"
FEATURE_EXTRACTION_BASE_URL = "https://api-inference.huggingface.co/pipeline/feature-extraction/"
FEATURE_EXTRACTION_MODEL_NAME = "sentence-transformers/paraphrase-MiniLM-L6-v2"

class HuggingFaceInferenceApiWrapper():
    def __init__(self, model_name: str, api_key: str) -> None:
        self.inference_base_url = INFERENCE_BASE_URL
        self.feature_extraction_base_url = FEATURE_EXTRACTION_BASE_URL
        self.model_name = model_name
        self.feature_extraction_model_name = FEATURE_EXTRACTION_MODEL_NAME
        self.api_key = api_key
        self.global_params = {
            "max_new_tokens": 512,
            "return_full_text": False,
            "repetition_penalty":1.1,
            "temperature": 0.01,
            "use_cache": False,
            "wait_for_model": True
        }
        self.global_options = {
            "use_cache": False,
            "wait_for_model": True
        }

    def predict(self, prompt: str) -> (dict, int):
        res = requests.post(
            url=f"{self.inference_base_url}{self.model_name}",
            headers={
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json",
                },
            json={
                "inputs": prompt,
                "parameters": self.global_params,
                "options": self.global_options
            }
        )
        parsed_response = res.json()
        try:
            processed_response = self._process_response(parsed_response[0]['generated_text'])
        except:
            processed_response = ""
        with open(PROMPT_LOG_PATH, "a") as f:
            date = datetime.datetime.now()
            f.write(f"[{date}][{version}]Prompt: {prompt}\nResponse: {processed_response}\n\n")
        return processed_response, res.status_code
    
    def get_embeddings(self, prompt: list) -> (list,int):
        res = requests.post(
            url=f"{self.feature_extraction_base_url}{self.feature_extraction_model_name}",
            headers={
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json",
                },
            json={
                "inputs": prompt,
            }
        )
        parsed_response = res.json()
        return parsed_response, res.status_code

    def _process_response(self, response: str) -> dict:
        try:
            response = response.replace("\n", "").replace("\t", "").replace("\\\"", "")
            response = response.replace("Assistant: ","")
            json_pattern = r'{.*}'
            matches = re.findall(json_pattern, response)
            return json.loads(matches[0])
        except json.decoder.JSONDecodeError or IndexError as e:
            print(response)
            print(e)
            return None