import os
import re
import uuid
from elasticsearch import Elasticsearch, RequestsHttpConnection
from elasticsearch_dsl import Search


class Db:

    def setup(self):
        self.es = Elasticsearch(
            hosts=[os.getenv("ELASTICSEARCH", "elasticsearch:9200")],
            use_ssl=True if os.getenv("ELASTICSEARCH_USE_SSL", "false").lower() == "true" else False,
            verify_certs=True if os.getenv("ELASTICSEARCH_VERIFY_CERTS", "false").lower() == "true" else False,
            connection_class=RequestsHttpConnection,
            http_auth=(os.getenv("ELASTICSEARCH_USER", 'user'), os.getenv("ELASTICSEARCH_PASS", 'pass')),
            # sniff_on_start=False,
            # sniff_on_connection_fail=False,
            # sniffer_timeout=1,
            # sniff_timeout=1,
            max_retries=1
        )

        # print(self.es.cluster.health())

    def parse_item(self, item):
        parsed_item = item["_source"]

        if "_score" in item:
            parsed_item["_score"] = item["_score"]

        if "_explanation" in item:
            parsed_item["_explanation"] = item["_explanation"]

        return parsed_item

    def highlight_string_in_field(self, item, filter, highlight_start="<em>", highlight_end="</em>"):
        p = re.compile("(%s)" % (filter), re.IGNORECASE)
        item = p.sub('%s\\1%s' % (highlight_start, highlight_end), item)

        return item

    def list(self, filter, sort, tags=[]):

        search_params = {
            "index": "notes",
            "size": 10000,
            "explain": True,
            "track_scores": True,
            "sort": ""
        }

        filter_query = []

        if len(tags) > 0:
            filter_query.append("(tags: (%s) )" % (" AND ".join(tags)))

        filter_string = []
        if len(filter) > 0:
            filter_terms = filter.split()

            # filter_string.append("(content.keyword: %s OR title.keyword: %s OR tags.keyword: %s)" % (filter, filter, filter))

            for term in filter_terms:
                filter_string.append("(content: *%s* OR title: *%s*)" % (term, term))
                filter_string.append("(content: %s OR title: %s)" % (term, term))

            filter_string = (" OR ").join(filter_string)

            filter_query.append(filter_string)

        filter_query = (" AND ").join(filter_query)

        if len(filter_query) > 0:
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
        res = self.es.get(index="notes", id=id)
        return self.parse_item(res)

    def create(self, data):
        data["id"] = str(uuid.uuid4())
        self.update(data)
        return data

    def update(self, data):
        self.es.index(index="notes", id=data["id"], body=data, refresh="wait_for")

    def delete(self, id):
        note = self.es.get(index="notes", id=id)
        self.es.delete(index="notes", id=id, refresh="wait_for")

        return self.parse_item(note)

    def tags(self):
        self.search = Search(using=self.es, index="notes")

        query = {
            "size": 0,
            "aggs": {
                "tags": {
                    "terms": {
                        "field": "tags.keyword",
                        "order": {
                            "_key": "asc"
                        },
                        "size": 500
                    }
                }
            }
        }

        self.search = Search.update_from_dict(self.search, query)
        result = self.search.execute()

        return list(map(lambda item: item.key, result.aggregations.tags))
