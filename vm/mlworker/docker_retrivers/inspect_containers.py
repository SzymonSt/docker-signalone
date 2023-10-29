import uuid
from docker import DockerClient
from mlworker.docker_retrivers.logs import retrive_logs
from common.issues_manager import IssuesManager

def inspect(dc: DockerClient, im: IssuesManager) -> list:
    parsed_containers = []
    issues = []
    containers = dc.containers.list(all=True)
    for container in containers:
        container_dict = container.__dict__
        del container_dict['client']
        del container_dict['collection']
        del container_dict['attrs']['ResolvConfPath']
        parsed_containers.append(container_dict)
        state = get_state(container_dict)
        if state['Status'] != 'running':
            if (state['Dead'] or
                state['OOMKilled'] or
                state['ExitCode'] != 0 or
                state['Error'] != '' or
                (state['Restarting'] and get_restart_count(container_dict) > 0)
            ):
                issue = {
                    "id": uuid.uuid1().hex,
                    "container_id": container_dict['attrs']['Id'],
                    "issue_type": "error",
                    "issue_severity": "critical",
                    "is_resolved": False,
                    "timestamp": container_dict['attrs']['State']['FinishedAt'],
                    "container_state": state['Error'],
                    "logs": retrive_logs(container_dict['attrs']['Id'], dc)
                }
                current_issue = im.get_issues_by_container_id(issue['container_id'])
                print(current_issue)
                if len(current_issue) == 0:
                    issues.append(issue)
                else:
                    if current_issue[0][5] != issue['timestamp']:
                        im.drop_issues(current_issue)
                        issues.append(issue)
    return issues


def get_state(container) -> dict:
    state = container['attrs']['State']
    del state['StartedAt']
    return state

def get_restart_count(container) -> int:
    restart_count = container['attrs']['RestartCount']
    return restart_count