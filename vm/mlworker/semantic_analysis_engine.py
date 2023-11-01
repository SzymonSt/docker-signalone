from mlworker.context_retrivers.db import vector_search
from qdrant_client import QdrantClient
from mlworker.huggingface_inference_api.wrapper import HuggingFaceInferenceApiWrapper
from mlworker.prompts.templates_llama import prompt_template_state, prompt_template_anomaly

class SemanticAnalysisEngine:
    def __init__(self, main_model_name, hf_api_key, qclient_url):
        self.qdrant_client = QdrantClient(
            qclient_url
        )
        self.hf_issues_inference_model = HuggingFaceInferenceApiWrapper(
            model_name=main_model_name,
            api_key=hf_api_key
        )

    def analyze(self, issues):
        for issue in issues:
            embeddings_prompt = ""
            if issue['container_state'] != "":
                embeddings_prompt += str(issue['container_state'])
            if issue['logs'] != "":
                embeddings_prompt +=  "\n" + str(issue['logs'])
            status = 0
            while status != 200:
                print("Retrieving embeddings...")
                ctxvec, status = self.hf_issues_inference_model.get_embeddings(embeddings_prompt)
            print("Retrival complete")
            context = vector_search(self.qdrant_client, ctxvec)
            status = 0
            if issue['issue_type'] == "error":
                while status != 200:
                    print("Predicting issue solutions...")
                    mess, status = self.hf_issues_inference_model.predict(prompt=prompt_template_state.format(
                        issue=str(issue), 
                        context=str(context[0].payload.get('possible_solutions'))))
                print("Prediction complete")
            elif issue['issue_type'] == "anomaly":
                while status != 200:
                    print("Predicting issue solutions...")
                    mess, status = self.hf_issues_inference_model.predict(prompt=prompt_template_anomaly.format(
                        issue=str(issue['container_state']),logs=str(issue['logs']), 
                        context=str(context[0].payload.get('possible_solutions'))))
                print("Prediction complete")
            del issue['logs']
            del issue['container_state']
            issue['issue'] = mess['issue']
            issue['solutions'] = mess['solutions']
            