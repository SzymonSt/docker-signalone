import json
from pymongo import MongoClient
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from dotenv import dotenv_values
import time
from dateutil import parser


def main():
    env_vars = dotenv_values(dotenv_path='.default.env')
    db_name = env_vars.get('APPLICATION_DB_NAME')
    coll_name = env_vars.get('APPLICATION_USERS_COLLECTION_NAME')
    issues_coll_name = env_vars.get('APPLICATION_ISSUES_COLLECTION_NAME')
    print(f"db_name: {db_name}, coll_name: {coll_name}, issues_coll_name: {issues_coll_name}")
    mongo_client = None

    while mongo_client is None:
        mongo_client = MongoClient('mongodb://localhost:27017/')
        time.sleep(1)

    
    with open('./devenv/users.json','r') as users_file:
        users = json.load(users_file)
    with open('./devenv/analysis.json','r') as issues_file:
        issues = json.load(issues_file)
    
    init_coll(users, mongo_client, db_name, coll_name)
    init_coll(issues, mongo_client, db_name, issues_coll_name)

def init_coll(objects, mongo_client, db_name, coll_name):
    try:
        if objects[0]["timestamp"] is not None:
            for obj in objects:
                obj["timestamp"] = parser.parse(obj["timestamp"])
    except KeyError:
        pass
    db = mongo_client[db_name]
    coll = db[coll_name]
    res = coll.insert_many(objects)
    check = coll.find({})
    if check is not None:
        print(f"Inserted {len(res.inserted_ids)} objects to {coll_name} collection")

if __name__ == '__main__':
    main()