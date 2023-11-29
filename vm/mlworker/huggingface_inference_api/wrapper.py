import requests

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
            "max_new_tokens": 1024,
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
                "inputs": prompt
            }
        )
        parsed_response = res.json()
        return parsed_response, res.status_code