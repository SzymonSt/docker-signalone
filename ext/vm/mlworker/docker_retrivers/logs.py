from docker import DockerClient

def retrive_logs(container_id: str, docker_client: DockerClient) -> str:
    container = docker_client.containers.get(container_id)
    logs = str(container.logs())
    return logs