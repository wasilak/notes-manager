import json
import os
from datetime import datetime
from flask import Flask, render_template, jsonify, request, send_from_directory
import importlib


# getting info about app version from package.json
with open("./package.json") as json_file:
    package_json = json.load(json_file)

app = Flask(__name__)

db_provider = os.getenv("DB_PROVIDER", "file")
db_module = importlib.import_module("library.db_providers.%s" % db_provider)
Db = db_module.Db
db = Db()


def connection():
    try:
        db.setup()
    except Exception as e:
        return str(e)

    return False


db_conn_err = connection()


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

    db.update(updated_note["id"], updated_note)

    return jsonify({"data": updated_note})


@app.route('/api/note/delete/<uuid>', methods=['DELETE'])
def api_note_delete(uuid):
    return jsonify({"data": db.delete(uuid)})


@app.route('/api/note/new', methods=['POST'])
def api_note_new():

    dt = datetime.today()
    seconds = int(dt.timestamp())

    new_note = json.loads(request.data)["note"]
    new_note["created"] = seconds
    new_note["updated"] = seconds

    new_note = db.create(new_note["id"], new_note)

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
