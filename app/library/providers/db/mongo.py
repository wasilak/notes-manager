import os
from bson.objectid import ObjectId
import pymongo
import logging
from urllib import parse


class Db:

    def setup(self):
        uri = "mongodb://%s:%s@%s" % (
            parse.quote_plus(os.getenv("MONGO_USER", 'user')),
            parse.quote_plus(os.getenv("MONGO_PASS", 'pass')),
            os.getenv("MONGO_HOST", "localhost:27017")
        )

        self.client = pymongo.MongoClient(uri)
        self.db = self.client.notes

        self.setup_indices()

        self.logger = logging.getLogger("uvicorn.error")

    def setup_indices(self):
        cur_indices = list(self.db.notes.list_indexes())

        text_index_name = "content_text_title_text"

        text_index_created = False

        for index in cur_indices:
            if index["name"] == text_index_name:
                text_index_created = True

        if not text_index_created:
            self.logger.info("Creating text index")
            self.db.notes.create_index([('content', 'text'), ('title', 'text')])


    def parse_item(self, doc):
        if "_id" in doc:
            doc["id"] = str(doc["_id"])
            del(doc["_id"])

        if "score" in doc:
            doc["_score"] = doc["score"]
            del(doc["score"])

        return doc

    def list(self, filter, sort, tags=[]):

        search_params = {}
        other_params = {}
        sort_params = []


        if len(tags) > 0:
            search_params["tags"] = {"$all": tags}

        if len(filter) > 0:
            search_params["$text"] = {"$search": filter}
            other_params["score"] = {"$meta": "textScore"}
            sort_params.append(('score', {'$meta': 'textScore'}))

        if len(sort) > 0:
            sort_tmp = sort.split(":")
            sort_order = 1 if sort_tmp[1] == "asc" else -1
            sort_params.append((sort_tmp[0], sort_order))

        if len(sort_params) == 0:
            # workaround to make searches without sort work
            sort_params.append(("$natural", pymongo.ASCENDING))

        docs = list(self.db.notes.find(search_params, other_params).sort(sort_params))

        docs = list(map(lambda doc: self.parse_item(doc), docs))

        return docs

    def get(self, id):
        doc = self.db.notes.find_one({"_id": ObjectId(id)})
        return self.parse_item(doc)

    def create(self, data):
        if "id" in data:
            del(data["id"])
        data["_id"] = ObjectId()
        self.db.notes.insert_one(data)
        return self.parse_item(data)

    def update(self, data):
        data["_id"] = ObjectId(data["id"])
        del(data["id"])
        self.db.notes.replace_one({"_id": data["_id"]}, data)

        self.parse_item(data)

    def delete(self, id):
        doc = self.get(id)
        self.db.notes.delete_one({"_id": ObjectId(id)})

        return self.parse_item(doc)

    def tags(self):
        return list(self.db.notes.distinct("tags"))
