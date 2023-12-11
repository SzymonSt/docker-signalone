import os

def updateMinorVersion(current_version: str):
    version = current_version.split('.')
    version[-2] = str(int(version[-2]) + 1)
    return '.'.join(version)

def getNewVersion(file_template, dir='../sources'):
    files = os.listdir(dir)
    files_matching_template = [file for file in files if file_template in file]
    files_matching_template.sort()
    if len(files_matching_template) == 0:
        return "v0.0.1"
    else:
        return files_matching_template[-1].split('_')[-1].replace('.csv', '')