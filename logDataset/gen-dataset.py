import csv
import json
from openai import OpenAI
from dotenv import load_dotenv
from tqdm import tqdm

# Load .env file
load_dotenv()

class DataGenerator:
    def __init__(self):
        self.prompt = "I have a problem with my Docker logs. Can you help me troubleshoot? Give solution for this anomaly from these docs and also give heading as the type of anamoly."
        self.temperature = 0.8
        self.max_tokens = 2048
        self.model_engine = "gpt-3.5-turbo"
        self.client = OpenAI()

    def generate_logs(self):
        prompt = """Your job is to create a synthetic dataset for training of bart model for log summarization of logs and their summaries. Create 5 log summary pair of synthetic logs for any technical anomaly and create their summaries in pairs. Create long logs of about 10-15 lines and then summarize them

                    sample
                    {"logs":"Logs of any software which are failing with error code, function name and error code"
                    "summary":
                    "Summarize the logs and what is happening in high detail. Include log details."}

                    only return a json format file and nothing else. Do not return any other text except the json. Do not give any starting or trailing text that is not data. Only generate 5 pair of logs and summaries. Do not include any commas inside the logs or summary. Always imagine the names of software or functions which give logs."""
        chat_completion = self.client.chat.completions.create(
                messages=[
                {
                    "role": "user",
                    "content": prompt,
                }
            ],
            model="gpt-3.5-turbo",
        )
       
        try:
            chat_completion = self.client.chat.completions.create(
                messages=[
                    {
                        "role": "user",
                        "content": prompt,
                    }
                ],
                model="gpt-3.5-turbo",
            )
            parsed_data = json.loads(str(chat_completion.choices[0].message.content),strict=False,)
            print(parsed_data)
            # Extract logs and summaries from data
            logs = [entry['logs'] for entry in parsed_data['data']]
            summary = [entry['summary'] for entry in parsed_data['data']]
        except Exception as e:
            print("Error:", e)
            return None, None
 
        
        return logs,summary


if __name__ == "__main__":
    # Create DataGenerator instance
    data_generator = DataGenerator()

    # Generate and save logs and summaries in a CSV file
    with open('dataset.csv', 'a', newline='') as csvfile:
        fieldnames = ['logs', 'summary']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        
        writer.writeheader()

        # Generate and write data 200 times. So, we will have 200 * 5 logs and summaries
        for _ in tqdm(range(200), desc="Generating data", unit=" logs"):
            logs, summary = data_generator.generate_logs()
            if logs is not None and summary is not None:
                for log, sum in zip(logs, summary):
                    writer.writerow({'logs': log, 'summary': sum})

    print("CSV file created successfully.")