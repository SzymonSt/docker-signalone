import uvicorn
import os
import json
from dotenv import load_dotenv
from fastapi import FastAPI, Response
from fastapi.middleware.cors import CORSMiddleware
from common.issues_manager import IssuesManager


def main():
    load_dotenv()
    db_path = os.getenv('DB_PATH')
    issues_manager = IssuesManager(db_path)
    app = FastAPI()
    app.add_middleware(CORSMiddleware, allow_origins=['*'])

    @app.get('/issues')
    async def analyze(
        container_id: str = None, 
        issue_type: str = None,
        issue_severity: str = None,
        is_resolved: str = None,
        start_timestamp: str = None,
        end_timestamp: str = None,
        offset: int = 0,
        limit: int = 100):
        parsed_issues = []
        filters = []
        if issue_type:
            filters.append({'key': 'issue_type', 'value': issue_type})
        if issue_severity:
            filters.append({'key': 'issue_severity', 'value': issue_severity})
        if is_resolved:
            filters.append({'key': 'is_resolved', 'value': is_resolved})
        if start_timestamp:
            filters.append({'key': 'timestamp', 'value': start_timestamp})
        if end_timestamp:
            filters.append({'key': 'end_timestamp', 'value': end_timestamp})

        if len(filters) == 0:
            filters = None

        issues = issues_manager.get_issues_by_container_id(
            container_id=container_id,
            filters=filters,
            offset=offset,
            limit=limit)
        
        for issue in issues:
            solutions = issue[7].replace("\"","'").replace("[","").replace("]","").split("', '")
            parsed_issue = {
                'id': issue[0],
                'container_id': issue[1],
                'issue_type': issue[2],
                'issue_severity': issue[3],
                'is_resolved': issue[4],
                'timestamp': issue[5],
                'issue': issue[6],
                'solutons': solutions
            }
            parsed_issues.append(parsed_issue)
        
        return Response(content=json.dumps(parsed_issues), media_type='application/json')

    @app.get('/containers')
    async def get_containers():
        parsed_containers = []
        containers = issues_manager.get_containers()
        for container in containers:
            parsed_containers.append(container[0])
        return Response(content=json.dumps(parsed_containers), media_type='application/json')

    uvicorn.run(app, host='0.0.0.0', port=8000)

if __name__ == '__main__':
    main()