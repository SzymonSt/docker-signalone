import logging
from mlworker.huggingface_inference_api.wrapper import HuggingFaceInferenceApiWrapper

class SemanticAnalysisEngine:
    def __init__(self, main_model_name, hf_api_key):
        self.hf_issues_inference_model = HuggingFaceInferenceApiWrapper(
            model_name=main_model_name,
            api_key=hf_api_key
        )

    def analyze(self, issues):
        for issue in issues:
            embeddings_prompt = "summarize: "
            if issue['logs'] != "":
                embeddings_prompt +=  str(issue['logs'])

            if len(embeddings_prompt) > 1024:
                embeddings_prompt = embeddings_prompt[-1024:]
            status = 0
            if issue['issue_type'] == "error":
                while status != 200:
                    print("Predicting issue solutions...")
                    logging.info("Predicting issue solutions...")
                    mess, status = self.hf_issues_inference_model.predict(prompt=embeddings_prompt)
                    if status != 200:
                        print("Retrival failed. Reason: {}".format(mess))
                        logging.info("Retrival failed. Reason: {}".format(mess))
                print("Prediction complete")
                logging.info("Prediction complete")
            elif issue['issue_type'] == "anomaly":
                while status != 200:
                    print("Predicting issue solutions...")
                    mess, status = self.hf_issues_inference_model.predict()
                print("Prediction complete")
                logging.info("Prediction complete")
            del issue['logs']
            del issue['container_state']
            issue['issue'] = 'TEST ISSUE'
            issue['solutions'] = mess[0]['summary_text']
            