import json
import os
from datetime import datetime
from flask import Flask, render_template, jsonify, request, send_from_directory
import re
import importlib


# getting info about app version from package.json
with open("./package.json") as json_file:
    package_json = json.load(json_file)

app = Flask(__name__)

db_provider = os.getenv("DB_PROVIDER", "file")
db_module = importlib.import_module("library.db_providers.%s" % db_provider)
Db = db_module.Db
db = Db()

storage_provider = os.getenv("STORAGE_PROVIDER", "none")
storage_module = importlib.import_module(
    "library.storage_providers.%s" % storage_provider)
Storage = storage_module.Storage
storage = Storage()


def connection():
    try:
        db.setup()
    except Exception as e:
        return str(e)

    return False


db_conn_err = connection()


def get_all_image_urls(content):
    pattern = re.compile(r'(https?:[/|.|\w|\s|-]*\.(jpg|gif|png|jpeg))')
    match = pattern.findall(content)

    match = list(set(match))

    result = []
    for url in match:
        result.append({
            "original": {
                "url": url[0],
                "extension": url[1]
            },
            "replacement": "",
        })

    return result


def replace_urls(content, image_urls):
    for item in image_urls:
        print(item["original"]["url"])
        print(item["replacement"])
        content = content.replace(item["original"]["url"], item["replacement"])

    return content


@app.route('/storage/<path:filename>')
def storage_uri(filename):
    return send_from_directory(app.root_path + '/storage/', filename)


@app.route('/static/node_modules/<path:filename>')
def base_static(filename):
    return send_from_directory(app.root_path + '/node_modules/', filename)


@app.route('/', defaults={'path': ''}, methods=['GET'])
@app.route('/<path:path>', methods=['GET'])
def index(path):
    return render_template('index.html', app_version=package_json["version"])


@app.route('/api/list/<filter>', methods=['GET'])
@app.route('/api/list/', defaults={'filter': ''}, methods=['GET'])
def api_list(filter):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return jsonify({"error": db_conn_err_persisting}), 503

    sort = request.args.get('sort', default="", type=str)

    if request.args.get('tags', type=str):
        tags = request.args.get('tags', type=str).strip().split(",")
    else:
        tags = []

    notes = db.list(filter.lower(), sort, tags)

    return jsonify({"data": notes})


@app.route('/api/note/<uuid>', methods=['GET'])
def api_note(uuid):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return jsonify({"error": db_conn_err_persisting}), 503

    note = db.get(uuid)

    return jsonify({"data": note})


@app.route('/api/note/<uuid>', methods=['POST'])
def api_note_update(uuid):

    dt = datetime.today()
    seconds = int(dt.timestamp())

    updated_note = json.loads(request.data)["note"]
    updated_note["updated"] = seconds

    if storage_provider != "none":
        image_urls = get_all_image_urls(updated_note["content"])
        storage.create_path(uuid)
        storage.get_files(uuid, image_urls)
        updated_note["content"] = replace_urls(updated_note["content"], image_urls)

    db.update(updated_note)

    return jsonify({"data": updated_note})


@app.route('/api/note/delete/<uuid>', methods=['DELETE'])
def api_note_delete(uuid):
    storage.cleanup(uuid)
    return jsonify({"data": db.delete(uuid)})


@app.route('/api/note/new', methods=['POST'])
def api_note_new():

    dt = datetime.today()
    seconds = int(dt.timestamp())

    new_note = json.loads(request.data)["note"]
    new_note["created"] = seconds
    new_note["updated"] = seconds

    new_note = db.create(new_note)

    if storage_provider != "none":
        # note first has to be created, in  order to have it's ID/_id
        # and afterwards images will have to be parsed and downloaded
        # and note itself - updated.
        image_urls = get_all_image_urls(new_note["content"])
        storage.create_path(new_note["id"])
        storage.get_files(new_note["id"], image_urls)
        new_note["content"] = replace_urls(new_note["content"], image_urls)
        db.update(new_note)

    return jsonify({"data": new_note})


@app.route('/api/tags/<query>', methods=['GET'])
@app.route('/api/tags/', defaults={'query': ''}, methods=['GET'])
def api_tags(query):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return jsonify({"error": db_conn_err_persisting}), 503

    # list(set()) removes duplicates from list
    tags = list(set(db.tags()))

    # filtering tags
    if len(query) > 0:
        tags = list(filter(lambda tag: query.lower() in tag.lower(), tags))

    tags.sort()

    return jsonify(tags)
