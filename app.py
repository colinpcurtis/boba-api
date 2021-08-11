from flask import Flask
from flask_cors import CORS, cross_origin
import json
from pymongo import MongoClient
import os

app = Flask("app")

mongoclient = MongoClient(os.getenv("MONGO_URL"))

db = mongoclient.boba_db

boba_count_db = db.boba_count

CORS(app)

app = Flask(__name__)


def get_users_in_server(collection):
    users = collection.find({})
    users_in_server = []
    for item in users:
        username = item["user"]
        users_in_server.append(username)
    return users_in_server


def does_server_exists(server):
    if server not in db.list_collection_names():
        return False
    return True


@app.route('/boba/<server>', methods=['GET'])
@cross_origin()
def get_boba(server: str):
    print(server)
    if not does_server_exists(server):
        return json.dumps({"error": f"Error: server {server} does not exist, please try again"})

    collection = db[server]
    users = get_users_in_server(collection)
    counts = {}
    cursor = boba_count_db.find({"user": {"$in": users}})
    # get the document in the count db
    # for all the users in the server
    for doc in cursor:
        user = doc["user"]
        count = doc["boba_count"]
        counts[user] = count
    count = {k: v for k, v in sorted(counts.items(), key=lambda item: item[1], reverse=True)}
    print(count)
    return json.dumps(counts)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000, debug=True)
