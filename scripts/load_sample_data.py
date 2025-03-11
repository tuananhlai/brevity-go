import pandas as pd
import requests
from elasticsearch import Elasticsearch
from datetime import datetime
import time
import io

# Initialize Elasticsearch client
es = Elasticsearch('http://localhost:9200')

def download_sample_data():
    """Download IMDB Movies sample data"""
    print("Downloading sample data...")
    url = "https://raw.githubusercontent.com/prust/wikipedia-movie-data/master/movies.json"
    response = requests.get(url)
    return response.json()

def create_index():
    """Create index with proper mapping"""
    print("Creating index...")
    index_name = "movies"
    
    # Delete index if it exists
    if es.indices.exists(index=index_name):
        es.indices.delete(index=index_name)
    
    # Define mapping
    mapping = {
        "mappings": {
            "properties": {
                "title": {"type": "text"},
                "year": {"type": "integer"},
                "cast": {"type": "keyword"},
                "genres": {"type": "keyword"},
                "directors": {"type": "keyword"}
            }
        }
    }
    
    # Create index with mapping
    es.indices.create(index=index_name, body=mapping)

def insert_data(data):
    """Insert data into Elasticsearch"""
    print("Inserting data...")
    index_name = "movies"
    batch_size = 1000
    
    # Process data in batches
    for i in range(0, len(data), batch_size):
        batch = data[i:i + batch_size]
        bulk_data = []
        
        for doc in batch:
            # Prepare bulk operation
            bulk_data.append({
                "index": {
                    "_index": index_name
                }
            })
            bulk_data.append(doc)
        
        # Execute bulk insert
        es.bulk(operations=bulk_data)
        print(f"Inserted {min(i + batch_size, len(data))} records...")
        
        # Small delay to prevent overwhelming Elasticsearch
        time.sleep(0.1)

def main():
    try:
        # Wait for Elasticsearch to be ready
        print("Waiting for Elasticsearch to be ready...")
        max_retries = 30
        retry_count = 0
        
        while retry_count < max_retries:
            try:
                if es.ping():
                    print("Successfully connected to Elasticsearch")
                    break
            except Exception:
                retry_count += 1
                print(f"Waiting for Elasticsearch... ({retry_count}/{max_retries})")
                time.sleep(2)
        
        if retry_count == max_retries:
            raise Exception("Could not connect to Elasticsearch after maximum retries")
        
        # Create index
        create_index()
        
        # Download and insert data
        data = download_sample_data()
        insert_data(data)
        
        # Verify data
        count = es.count(index="movies")
        print(f"\nSuccessfully inserted {count['count']} records into Elasticsearch!")
        
    except Exception as e:
        print(f"An error occurred: {str(e)}")

if __name__ == "__main__":
    main() 