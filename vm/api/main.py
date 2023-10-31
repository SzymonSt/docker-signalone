import uvicorn
import os
import json
from dotenv import load_dotenv
from fastapi import FastAPI
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware
from common.issues_manager import IssuesManager
class Issue(BaseModel):
    id: str
    containerId: str
    issueType: str
    issueSeverity : str
    isResolved: bool
    timestamp: str
    issue: str
    solutions: list[str]

def main():
    load_dotenv()
    db_path = os.getenv('DB_PATH')
    issues_manager = IssuesManager(db_path)
    app = FastAPI()
    app.add_middleware(CORSMiddleware,
        allow_origins=['*'],
        allow_credentials=True,
        allow_methods=['*'],)

    @app.get('/issues', response_model=list[Issue])
    async def analyze(
        containerId: str = None, 
        issueType: str = None,
        issueSeverity: str = None,
        isResolved: str = None,
        startTimestamp: str = None,
        endTimestamp: str = None,
        offset: int = 0,
        limit: int = 100):
        parsed_issues = []
        filters = []
        if issueType:
            filters.append({'key': 'issue_type', 'value': issueType})
        if issueSeverity:
            filters.append({'key': 'issue_severity', 'value': issueSeverity})
        if isResolved:
            filters.append({'key': 'is_resolved', 'value': isResolved})
        if startTimestamp:
            filters.append({'key': 'timestamp', 'value': startTimestamp})
        if endTimestamp:
            filters.append({'key': 'end_timestamp', 'value': endTimestamp})

        if len(filters) == 0:
            filters = None

        issues = issues_manager.get_issues_by_container_id(
            container_id=containerId,
            filters=filters,
            offset=offset,
            limit=limit)
        
        for issue in issues:
            solutions = issue[7].replace("\"","'").replace("[","").replace("]","").split("', '")
            parsed_issue = Issue(
                id=issue[0],
                containerId=issue[1],
                issueType=str.capitalize(issue[2]),
                issueSeverity=str.capitalize(issue[3]),
                isResolved=bool(issue[4]),
                timestamp=issue[5],
                issue=issue[6],
                solutions=solutions)
            parsed_issues.append(parsed_issue)
        
        return parsed_issues

    @app.get('/containers', response_model=list[str])
    async def get_containers():
        parsed_containers = []
        containers = issues_manager.get_containers()
        for container in containers:
            parsed_containers.append(container[0])
        return parsed_containers

    uvicorn.run(app, host='0.0.0.0', port=8000)

if __name__ == '__main__':
    main()