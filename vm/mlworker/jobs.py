import datetime
from docker import DockerClient
import numpy as np
import uuid
import os
import logging
from dotenv import load_dotenv
from sklearn.ensemble import IsolationForest
from mlworker.docker_retrivers.inspect_containers import inspect
from mlworker.docker_retrivers.runtime import retrive_runtime
from mlworker.semantic_analysis_engine import SemanticAnalysisEngine
from common.issues_manager import IssuesManager


baselines = {}
parsed_stats = {}
load_dotenv()
hf_api_key = os.getenv("HF_API_KEY")
qclient_url = os.getenv("QCLIENT_URL")
main_model_name = os.getenv("HF_MAIN_MODEL_NAME")
sanalysis = SemanticAnalysisEngine(
    main_model_name=main_model_name,
    hf_api_key= hf_api_key,
    qclient_url=qclient_url
)
im = IssuesManager()
logging.basicConfig(filename='mlworker.log',level=logging.DEBUG)

def resource_usage_anomaly():
    print("Running resource usage anomaly detection")
    logging.info("Running resource usage anomaly detection")
    try:
        dc = DockerClient('unix://var/run/docker.sock')
    except Exception as e:
        logging.error("Docker client connection error: {}".format(e))
    running_containers = dc.containers.list(all=False)
    for container in running_containers:
        issues = []
        container_dict = container.__dict__
        print("Checking container: ", container_dict['attrs']['Id'])
        logging.info("Checking container: ", container_dict['attrs']['Id'])
        cpu_anomaly = False
        memory_anomaly = False
        sample_size = 15
        window_size = 900
        sample_index = 0
        #check if container has previous measurements
        if container_dict['attrs']['Id'] not in parsed_stats:
            parsed_stats[container_dict['attrs']['Id']] = {
                "cpu_usage": [],
                "memory_usage": []
            }
        #check max of sliding window
        if len(parsed_stats[container_dict['attrs']['Id']]['cpu_usage']) >= window_size:
            parsed_stats[container_dict['attrs']['Id']]['cpu_usage'] = parsed_stats[container_dict['attrs']['Id']]['cpu_usage'][30:]
        if len(parsed_stats[container_dict['attrs']['Id']]['memory_usage']) >= window_size:
            parsed_stats[container_dict['attrs']['Id']]['memory_usage'] = parsed_stats[container_dict['attrs']['Id']]['memory_usage'][30:]
        #new measurements
        for stat in container.stats(decode=True, stream=True):
            if sample_index > sample_size:
                break
            parsed_stats[container_dict['attrs']['Id']]['cpu_usage'].append(stat['cpu_stats']['cpu_usage']['total_usage'])
            parsed_stats[container_dict['attrs']['Id']]['memory_usage'].append(stat['memory_stats']['usage'])
            sample_index += 1
        #get baseline, create if not exists
        try:
            baseline = baselines[container_dict['attrs']['Id']]
            cpu_anomaly = baseline['cpu_anomaly']
            memory_anomaly = baseline['memory_anomaly']
        except:
            baseline = {
                "cpu": np.percentile(parsed_stats[container_dict['attrs']['Id']]['cpu_usage'],25),
                "cpu_anomaly": cpu_anomaly,
                "memory": np.percentile(parsed_stats[container_dict['attrs']['Id']]['memory_usage'],25),
                "memory_anomaly": memory_anomaly
            }
            baselines[container_dict['attrs']['Id']] = baseline
        print("Baseline: ", baseline)

        #CPU
        parsed_stats_cpu = parsed_stats[container_dict['attrs']['Id']]['cpu_usage']
        if not cpu_anomaly:   
            cpu_anomaly = _check_for_anomaly(parsed_stats_cpu, baseline, "cpu")
            if cpu_anomaly:
                logs = container.logs(tail=15)
                runtime = retrive_runtime(container_dict['attrs']['Config']['Cmd'], container_dict['attrs']['Config']['Env'])
                issue = {
                    "id": uuid.uuid1().hex,
                    "container_id": container_dict['attrs']['Id'],
                    "issue_type": "anomaly",
                    "issue_severity": "warning",
                    "is_resolved": False,
                    "timestamp": datetime.datetime.now(),
                    "container_state": """
                        Docker CPU resource excessive usage for runtime {runtime}, 
                        analyze provided logs of application to find functions or endpoints as possible root cause""".format(runtime=runtime),
                    "logs": logs
                }
                issues.append(issue)
                print("CPU anomaly detected")
            else:
                baseline['cpu'] = np.percentile(parsed_stats_cpu,25)
        else:
            if np.percentile(parsed_stats_cpu,25) <= baseline['cpu']:
                _recover_from_anomaly(container_dict['attrs']['Id'])
                cpu_anomaly = False 

        #Memory
        parsed_stats_mem = parsed_stats[container_dict['attrs']['Id']]['memory_usage']
        if not memory_anomaly:
            memory_anomaly = _check_for_anomaly(parsed_stats_mem, baseline, "memory")
            if memory_anomaly:
                logs = container.logs(tail=15)
                runtime = retrive_runtime(container_dict['attrs']['Config']['Cmd'], container_dict['attrs']['Config']['Env'])
                issue = {
                    "id": uuid.uuid1().hex,
                    "container_id": container_dict['attrs']['Id'],
                    "issue_type": "anomaly",
                    "issue_severity": "warning",
                    "is_resolved": False,
                    "timestamp": datetime.datetime.now(),
                    "container_state": """
                        Docker Memory resource excessive usage for runtime {runtime}, 
                        analyze provided logs of application to find functions or endpoints as possible root cause""".format(runtime=runtime),
                    "logs": logs
                }
                issues.append(issue)
                print("Memory anomaly detected")
            else:
                baseline['memory'] = np.percentile(parsed_stats_mem,25)
        else:
            if np.percentile(parsed_stats_mem,25) <= baseline['memory']:
                _recover_from_anomaly(container_dict['attrs']['Id'])
                memory_anomaly = False

        baseline['cpu_anomaly'] = memory_anomaly
        baseline['memory_anomaly'] = memory_anomaly
        sanalysis.analyze(issues)
        im.insert_issues(issues)
        baselines[container_dict['attrs']['Id']] = baseline

def container_error_scan():
    print("Running container error scan")
    logging.info("Running container error scan")
    dc = DockerClient().from_env()
    issues = inspect(dc, im)
    sanalysis.analyze(issues)
    im.insert_issues(issues)

def _calculate_outliners(data: list):
    data = np.array(data).reshape(-1, 1)
    clf = IsolationForest(random_state=0)
    return clf.fit_predict(data)

def _check_for_anomaly(parsed_stats: list, baseline: object, resource: str) -> bool:
    anomaly = False
    anomaly_res = _calculate_outliners(parsed_stats)
    outliners = []
    inliners = []
    for idx, measurement in enumerate(anomaly_res):
        if measurement == -1  and parsed_stats[idx] > baseline[resource]:
            outliners.append(parsed_stats[idx])
        else:
            inliners.append(parsed_stats[idx])
    outliners_avg = sum(outliners)/(len(outliners) if len(outliners) > 0 else 1)
    print("Outliners avg: ", outliners_avg)
    print("Outliners fraction:",len(outliners)/len(parsed_stats))
    if (sum(inliners)/len(inliners) if sum(inliners)/len(inliners) > 0 else 1)*4 < outliners_avg and len(outliners)/len(parsed_stats) >= 0.1:
        anomaly = True
    print("Anomaly detected: ", anomaly)
    return anomaly

def _recover_from_anomaly(container_id: str):
    filters=[]
    filters.append({
        "key": "issue_type",
        "value": "anomaly"
    })
    issues = im.get_issues_by_container_id(container_id, filters)
    im.drop_issues(issues)