import uuid
import json
import os
from datetime import datetime
from flask import Flask, render_template, jsonify, request, send_from_directory


app = Flask(__name__)

db_provider = os.getenv("DB_PROVIDER", "file")

if db_provider == "elasticsearch":
    from db_providers.elastic import Db

elif db_provider == "file":
    from db_providers.file import Db

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
    return render_template('index.html')


@app.route('/api/list/<filter>', methods=['GET'])
@app.route('/api/list/', defaults={'filter': ''}, methods=['GET'])
def api_list(filter):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return jsonify({"error": db_conn_err_persisting}), 503

    sort = request.args.get('sort', default="", type=str)

    notes = db.list(filter.lower(), sort)

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
    new_note["id"] = str(uuid.uuid4())
    new_note["created"] = seconds
    new_note["updated"] = seconds

    db.create(new_note["id"], new_note)

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
