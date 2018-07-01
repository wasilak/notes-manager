import uuid
import json
import time
import re
import os
from datetime import datetime

from flask import Flask, render_template, jsonify, request
from elasticsearch import Elasticsearch


app = Flask(__name__)

es = Elasticsearch(hosts=[os.getenv("ELASTICSEARCH", "elasticsearch:9200")])

try:
    es_health = es.cluster.health()
except Exception as e:
    print(e)
    exit(1)

print(es_health)


def parse_item(item):
    parsed_item = item["_source"]

    if "_score" in item:
        parsed_item["_score"] = item["_score"]

    if "_explanation" in item:
        parsed_item["_explanation"] = item["_explanation"]

    return parsed_item


def highlight_string_in_field(item, filter, highlight_start="<em>", highlight_end="</em>"):
    p = re.compile("(%s)" % (filter), re.IGNORECASE)
    print(item)
    item = p.sub('%s\\1%s' % (highlight_start, highlight_end), item)
    print(item)

    return item


@app.route('/', defaults={'path': ''}, methods=['GET'])
@app.route('/<path:path>', methods=['GET'])
def index(path):
    return render_template('index.html')


@app.route('/api/list/<filter>', methods=['GET'])
@app.route('/api/list/', defaults={'filter': ''}, methods=['GET'])
def api_list(filter):

    if len(filter) > 0:
        filter_terms = filter.split()

        filter_query = []
        for term in filter_terms:
            filter_query.append("(content: *%s* OR title: *%s*)" % (term, term))

        filter_query = (" OR ").join(filter_query)
        res = es.search(index="notes", doc_type='doc', size=10000, q=filter_query, explain=True)
    else:
        res = es.search(index="notes", doc_type='doc', size=10000, explain=True)

    parsed_items = []

    for item in res["hits"]["hits"]:
        parsed_item = parse_item(item)

        # if len(filter) > 0:
        #     parsed_item["content"] = highlight_string_in_field(parsed_item["content"], filter, "*", "*")
        #     parsed_item["title"] = highlight_string_in_field(parsed_item["title"], filter)

        parsed_items.append(parsed_item)

    return jsonify({"data": parsed_items})


@app.route('/api/note/<uuid>', methods=['GET'])
def api_note(uuid):
    res = es.get(index="notes", doc_type='doc', id=uuid)
    return jsonify({"data": parse_item(res)})


@app.route('/api/note/<uuid>', methods=['POST'])
def api_note_update(uuid):

    dt = datetime.today()
    seconds = int(dt.timestamp())

    updated_note = json.loads(request.data)["note"]
    updated_note["updated"] = seconds

    del updated_note["edit"]

    es.index(index="notes", doc_type='doc', id=updated_note["id"], body=updated_note)
    time.sleep(1)

    return jsonify({"data": updated_note})


@app.route('/api/note/delete/<uuid>', methods=['DELETE'])
def api_note_delete(uuid):
    note = es.get(index="notes", doc_type='doc', id=uuid)
    es.delete(index="notes", doc_type='doc', id=uuid)
    time.sleep(1)
    return jsonify({"data": parse_item(note)})


@app.route('/api/note/new', methods=['POST'])
def api_note_new():

    dt = datetime.today()
    seconds = int(dt.timestamp())

    new_note = json.loads(request.data)["note"]
    new_note["id"] = str(uuid.uuid4())
    new_note["created"] = seconds
    new_note["updated"] = seconds

    es.index(index="notes", doc_type='doc', id=new_note["id"], body=new_note)
    time.sleep(1)

    return jsonify({"data": new_note})
