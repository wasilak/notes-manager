import uuid
import json
import os
from datetime import datetime
from Crypto.Cipher import AES
import base64

from flask import Flask, render_template, jsonify, request, send_from_directory


app = Flask(__name__)

db_provider = os.getenv("DB_PROVIDER", "file")

secret_key = os.getenv("CRYPTO_SECRET_KEY", False)

if secret_key:
    cipher = AES.new(secret_key, AES.MODE_ECB)


def string_padding(text):
    extra = len(text) % 16
    if extra > 0:
        text = text + (" " * (16 - extra))

    return text


def string_encode(text):
    if secret_key:
        try:
            return base64.b64encode(cipher.encrypt(string_padding(text))).strip().decode('utf8')
        except Exception as e:
            print(e)

    return text


def string_decode(text):
    if secret_key:
        try:
            return cipher.decrypt(base64.b64decode(text)).strip().decode('utf8')
        except Exception as e:
            print(e)

    return text


def note_encode(note):
    note["encoded"] = True
    note["title"] = string_encode(note["title"])
    note["content"] = string_encode(note["content"])

    if "tags" in note:
        note["tags"] = list(map(lambda tag: string_encode(tag), note["tags"]))

    return note


def note_decode(note):
    note["title"] = string_decode(note["title"])
    note["content"] = string_decode(note["content"])

    if "tags" in note:
        note["tags"] = list(map(lambda tag: string_decode(tag), note["tags"]))

    return note


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

    for note in notes:
        note = note_decode(note)

    return jsonify({"data": notes})


@app.route('/api/note/<uuid>', methods=['GET'])
def api_note(uuid):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return jsonify({"error": db_conn_err_persisting}), 503

    note = note_decode(db.get(uuid))

    return jsonify({"data": note})


@app.route('/api/note/<uuid>', methods=['POST'])
def api_note_update(uuid):

    dt = datetime.today()
    seconds = int(dt.timestamp())

    updated_note = json.loads(request.data)["note"]
    updated_note["updated"] = seconds
    updated_note = note_encode(updated_note)

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

    new_note = note_encode(new_note)

    db.create(new_note["id"], new_note)

    return jsonify({"data": new_note})


@app.route('/api/tags', methods=['GET'])
def api_tags():

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return jsonify({"error": db_conn_err_persisting}), 503

    # list(set()) removes duplicates from list
    tags = list(set(list(map(lambda tag: string_decode(tag), db.tags()))))

    tags.sort()

    return jsonify(tags)
