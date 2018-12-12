import os
import re
from elasticsearch import Elasticsearch


class Db:

    def setup(self):
        self.es = Elasticsearch(
            hosts=[os.getenv("ELASTICSEARCH", "elasticsearch:9200")],
            # sniff_on_start=False,
            # sniff_on_connection_fail=False,
            # sniffer_timeout=1,
            # sniff_timeout=1,
            max_retries=1
        )
        print(self.es.cluster.health())

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

    def list(self, filter, sort):

        search_params = {
            "index": "notes",
            "doc_type": 'doc',
            "size": 10000,
            "explain": True,
            "track_scores": True,
            "sort": ""
        }

        if len(filter) > 0:
            filter_terms = filter.split()

            filter_query = []
            for term in filter_terms:
                filter_query.append("(content: *%s* OR title: *%s*)" % (term, term))

            filter_query = (" OR ").join(filter_query)

            search_params["q"] = filter_query

        if len(sort) > 0:
            search_params["sort"] = sort

        res = self.es.search(**search_params)

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
        self.es.index(index="notes", doc_type='doc', id=id, body=data, refresh="wait_for")

    def delete(self, id):
        note = self.es.get(index="notes", doc_type='doc', id=id)
        self.es.delete(index="notes", doc_type='doc', id=id, refresh="wait_for")

        return self.parse_item(note)
