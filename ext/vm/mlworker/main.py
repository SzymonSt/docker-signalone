import schedule
import time
import sqlite3

from mlworker.jobs import resource_usage_anomaly, container_error_scan

def main():
    # init
    db = sqlite3.connect("/data/issues.db")
    db.execute("""CREATE TABLE IF NOT EXISTS issues (
                id TEXT,
                container_id TEXT,
                issue_type TEXT,
                issue_severity TEXT,
                is_resolved BOOLEAN,
                timestamp TEXT,
                issue TEXT,
                solutions TEXT)""")
    db.commit()
    db.close()

    # Anomaly detection disabled for now to embed log summarization model 
    # schedule.every(5).seconds.do(resource_usage_anomaly)
    schedule.every(5).seconds.do(container_error_scan)                                    

    while True:
        schedule.run_pending()
        time.sleep(1)

if __name__ == '__main__':
    main()