import json
from pymongo import MongoClient
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from dotenv import dotenv_values
import time


def main():
    env_vars = dotenv_values(dotenv_path='.default.env')
    db_name = env_vars.get('APPLICATION_DB_NAME')
    coll_name = env_vars.get('APPLICATION_USERS_COLLECTION_NAME')
    sources_coll_name = env_vars.get('SOLUTION_COLLECTION_NAME')
    print(f"db_name: {db_name}, coll_name: {coll_name}, sources_coll_name: {sources_coll_name}")
    mongo_client = None
    qdrant_client = None

    while mongo_client is None or qdrant_client is None:
        mongo_client = MongoClient('mongodb://localhost:27017/')
        qdrant_client = QdrantClient(host='localhost', port=6333)
        time.sleep(1)

    
    with open('./devenv/users.json','r') as users_file:
        users = json.load(users_file)
    # with open('./devenv/solutions_sources.json','r') as solution_sources_file:
    #     solution_sources = json.load(solution_sources_file)
    
    init_users(users, mongo_client, db_name, coll_name)
    # init_solution_sources(solution_sources, qdrant_client, sources_coll_name)

def init_users(users, mongo_client, db_name, coll_name):
    db = mongo_client[db_name]
    coll = db[coll_name]
    res = coll.insert_many(users)
    check = coll.find({})
    if check is not None:
        print(f"Inserted {len(res.inserted_ids)} users to {coll_name} collection")
        for user in check:
            print(user)

def init_solution_sources(sources, qdrant_client, sources_coll_name):
    qdrant_client.recreate_collection(
        collection_name=sources_coll_name,
        vectors_config=VectorParams(size=384, distance=Distance.COSINE),
    )
    for si, source in enumerate(sources):
        vector = source['vector']
        payload = {
            'title': source['title'],
            'description': source['description'],
            'url': source['url'],
            'featuredAnswer': source['featuredAnswer'],
        }
        id = int(source['id'])
        try:
            res = qdrant_client.upsert(
                collection_name=sources_coll_name,
                points=[PointStruct(id=id, vector=vector, payload=payload)]
            )
            print(res)
            print(f"Upserted {si+1} out of {len(sources)} sources")
        except Exception as e:
            print(f"Failed to upsert source with id {id}")
            print(e)
            continue
    print(f"Inserted sources to {sources_coll_name} collection")

if __name__ == '__main__':
    main()