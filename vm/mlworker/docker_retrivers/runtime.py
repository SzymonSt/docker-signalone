from enum import Enum

runtime_dict = {
    "GOLANG": ["go", "GOLANG_VERSION", "GOPATH"],
    "PYTHON": ["python", "PYTHON_VERSION"],
}

def retrive_runtime(container_cmd: list[str], container_env: list[str]) -> str:
    container_runtime_props = []
    container_env = [env.split("=")[0] for env in container_env]
    container_runtime_props.extend(container_cmd)
    container_runtime_props.extend(container_env)
    for runtime in runtime_dict:
        if all(prop in container_runtime_props for prop in runtime_dict[runtime]):
            return runtime