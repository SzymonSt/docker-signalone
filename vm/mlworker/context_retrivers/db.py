
def vector_search(client, vector):
    hits = client.search(
    collection_name="stackoverflow_docker",
    query_vector=vector,
    limit=1
    )
    return hits