import json
from pymongo import MongoClient
from qdrant_client import QdrantClient, VectorParams, Distance, Point
from dotenv import dotenv_values


def main():
    env_vars = dotenv_values(dotenv_path='../.default.env')
    db_name = env_vars.get('APPLICATION_DB_NAME')
    coll_name = env_vars.get('APPLICATION_USERS_COLLECTION_NAME')
    mongo_client = MongoClient('mongodb://localhost:27017/')
    sources_coll_name = env_vars.get('SOLUTION_COLLECTION_NAME')
    qdrant_client = QdrantClient(host='localhost', port=6333)

    
    with open('./users.json','r') as users_file:
        users = json.load(users_file)
    with open('./solutions_sources.json','r') as solution_sources_file:
        solution_sources = json.load(solution_sources_file)
    
    init_users(users, mongo_client, db_name, coll_name)
    init_solution_sources(solution_sources, qdrant_client, sources_coll_name)

def init_users(users, mongo_client, db_name, coll_name):
    db = mongo_client[db_name]
    coll = db[coll_name]
    res = coll.insert_many(users)
    print(f"Inserted {len(res.inserted_ids)} users")

def init_solution_sources(sources, qdrant_client, sources_coll_name):
    qdrant_client.recreate_collection(
        collection_name=sources_coll_name,
        vectors_config=VectorParams(size=384, distance=Distance.COSINE),
    )
    for source in sources:
        vector = source['vector']
        payload = {
            'title': source['title'],
            'description': source['description'],
            'url': source['url'],
            'featuredAnswer': source['featuredAnswer'],
        }
        id = source['id']
        qdrant_client.upsert(
            collection_name=sources_coll_name,
            points=[Point(id=id, vector=vector, payload=payload)]
        )
    print(f"Inserted sources to {sources_coll_name} collection")

if __name__ == '__main__':
    main()