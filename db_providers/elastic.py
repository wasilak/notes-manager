import os
import time
import re
from elasticsearch import Elasticsearch


class Db:

    def __init__(self):
        try:
            self.es = Elasticsearch(hosts=[os.getenv("ELASTICSEARCH", "elasticsearch:9200")])
            print(self.es.cluster.health())
        except Exception as e:
            print(e)
            exit(1)

    def parse_item(self, item):
        parsed_item = item["_source"]

        if "_score" in item:
            parsed_item["_score"] = item["_score"]

        if "_explanation" in item:
            parsed_item["_explanation"] = item["_explanation"]

        return parsed_item

    def highlight_string_in_field(self, item, filter, highlight_start="<em>", highlight_end="</em>"):
        p = re.compile("(%s)" % (filter), re.IGNORECASE)
        print(item)
        item = p.sub('%s\\1%s' % (highlight_start, highlight_end), item)
        print(item)

        return item

    def list(self, filter):
        if len(filter) > 0:
            filter_terms = filter.split()

            filter_query = []
            for term in filter_terms:
                filter_query.append("(content: *%s* OR title: *%s*)" % (term, term))

            filter_query = (" OR ").join(filter_query)
            res = self.es.search(index="notes", doc_type='doc', size=10000, q=filter_query, explain=True)
        else:
            res = self.es.search(index="notes", doc_type='doc', size=10000, explain=True)

        parsed_items = []

        for item in res["hits"]["hits"]:
            parsed_item = self.parse_item(item)

            # if len(filter) > 0:
            #     parsed_item["content"] = highlight_string_in_field(parsed_item["content"], filter, "*", "*")
            #     parsed_item["title"] = highlight_string_in_field(parsed_item["title"], filter)

            parsed_items.append(parsed_item)

        return parsed_items

    def get(self, id):
        res = self.es.get(index="notes", doc_type='doc', id=id)
        return self.parse_item(res)

    def create(self, id, data):
        self.update(id, data)

    def update(self, id, data):
        self.es.index(index="notes", doc_type='doc', id=id, body=data)
        time.sleep(1)

    def delete(self, id):
        note = self.es.get(index="notes", doc_type='doc', id=id)
        self.es.delete(index="notes", doc_type='doc', id=id)
        time.sleep(1)

        return self.parse_item(note)
