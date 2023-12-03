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
        self.temperature = 0.9
        self.max_tokens = 2048
        self.model_engine = "gpt-3.5-turbo"
        self.client = OpenAI()

    def generate_logs(self):
        prompt = """Your job is to create a synthetic dataset for training of bart model for log summarization of logs and their summaries. Create 1 log summary pair of long synthetic logs for docker containers and create their summaries in pairs. Create long logs of about 10-15 lines and then summarize them. Do not use the sample in your results.

                    sample example
                    data[
                    {"logs":"2023-11-24T12:00:00Z [INFO] Server started on port 80
2023-11-24T12:05:00Z [ERROR] Internal Server Error: Database connection failed
2023-11-24T12:10:00Z [WARNING] Request from suspicious IP address: 192.168.1.100
2023-11-24T12:15:00Z [INFO] Request handled successfully
2023-11-24T12:20:00Z [ERROR] Out of memory: container restarting
2023-11-24T12:25:00Z [INFO] Server shutting down

2023-11-24T12:30:00Z [DEBUG] Database query executed in 50ms
2023-11-24T12:35:00Z [ERROR] Unauthorized access attempt: User 'admin' with incorrect password
2023-11-24T12:40:00Z [INFO] Container backup initiated
2023-11-24T12:45:00Z [WARNING] High CPU usage detected
2023-11-24T12:50:00Z [ERROR] File not found: /var/www/html/index.html
2023-11-24T12:55:00Z [INFO] Security patch applied successfully"
                    "summary":
                    "The Docker logs provide a chronological account of server activity. It begins with the successful startup of the server at 12:00:00, but issues arise at 12:05:00 with an internal server error due to a failed database connection. Subsequent logs at 12:10:00 note a warning for a suspicious IP address, and at 12:20:00, an out-of-memory error prompts the container to restart. The server initiates shutdown at 12:25:00. Additional logs include a database query execution time at 12:30:00, an unauthorized access attempt at 12:35:00, a container backup initiation at 12:40:00, a warning for high CPU usage at 12:45:00, a file not found error at 12:50:00, and the successful application of a security patch at 12:55:00. These logs collectively capture a range of events, including errors, warnings, and informational messages, providing insights into the server's performance and security aspects."}]

                    only return a json format file and nothing else. Do not return any other text except the json. Do not give any starting or trailing text that is not data. Only generate 1 pair of logs and summaries. Do not include any commas inside the logs or summary. Always imagine the names of software or functions which give logs."""
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

        # Generate and write data multiple times. So, we will have X * 5 logs and summaries
        for _ in tqdm(range(500), desc="Generating data", unit=" logs"):
            logs, summary = data_generator.generate_logs()
            if logs is not None and summary is not None:
                for log, sum in zip(logs, summary):
                    writer.writerow({'logs': log, 'summary': sum})

    print("CSV file created successfully.")  