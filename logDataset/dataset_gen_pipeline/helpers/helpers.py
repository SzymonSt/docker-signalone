

def updateMinorVersion(current_version: str):
    version = current_version.split('.')
    version[-2] = str(int(version[-2]) + 1)
    return '.'.join(version)