import json
import os
from datetime import datetime
from fastapi import FastAPI, BackgroundTasks
from fastapi.logger import logger
import re
import importlib

from starlette.requests import Request
from starlette.responses import Response, RedirectResponse
from dotenv import load_dotenv
from starlette.staticfiles import StaticFiles
from starlette.templating import Jinja2Templates
from pydantic import BaseModel
from typing import List

load_dotenv()

# getting info about app version from package.json
with open("./package.json") as json_file:
    package_json = json.load(json_file)

app = FastAPI(
    title="Notes Manager",
    description=package_json["description"],
    version=package_json["version"],
)

app.mount("/static", StaticFiles(directory="static"), name="static")

templates = Jinja2Templates(directory="templates")


class Note(BaseModel):
    id: str = None
    content: str
    title: str
    created: int = None
    updated: int = None
    _score: int = None
    tags: List[str] = []


db_provider = os.getenv("DB_PROVIDER", "file")
db_module = importlib.import_module("library.providers.db.%s" % db_provider)
Db = db_module.Db
db = Db()

storage_provider = os.getenv("STORAGE_PROVIDER", "none")
storage_module = importlib.import_module("library.providers.storage.%s" % storage_provider)
Storage = storage_module.Storage
storage = Storage()

if "local" == storage_provider:
    app.mount("/storage", StaticFiles(directory="storage"), name="storage")
if storage_provider in ["s3", "s3_minio"]:
    @app.get("/storage/{path:path}")
    async def storage_endpoint(request: Request, path: str = ''):
        presigned_url = storage.get_object(path)
        return RedirectResponse(url=presigned_url)


def connection():
    try:
        db.setup()
    except Exception as e:
        print(e)
        return str(e)

    return False


db_conn_err = connection()


def get_all_image_urls(content):
    pattern = re.compile(r'(https?:[\/\.\w\s\-\*]*\.(jpg|gif|png|jpeg))')
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
        content = content.replace(item["original"]["url"], item["replacement"])

    return content


@app.get("/api/list/")
async def api_list(tags: str = '', filter: str = '', sort: str = ''):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            return {"error": db_conn_err_persisting}

    cur_tags = []

    if len(tags) > 0:
        cur_tags = tags.strip().split(",")

    notes = db.list(filter.lower(), sort, cur_tags)

    return notes


@app.get("/api/note/{uuid}", response_model=Note)
async def api_note(uuid: str, response: Response):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            response.status_code = 503
            return {"error": db_conn_err_persisting}

    note = db.get(uuid)

    return note


def storage_get_files(note):
    image_urls = get_all_image_urls(note["content"])
    storage.get_files(note["id"], image_urls)
    note["content"] = replace_urls(note["content"], image_urls)
    db.update(note)


@app.post("/api/note/{uuid}", response_model=Note)
async def api_note_update(uuid: str, background_tasks: BackgroundTasks, item: Note):

    note = item.dict()

    dt = datetime.today()
    seconds = int(dt.timestamp())

    note["updated"] = seconds

    if storage_provider != "none":
        background_tasks.add_task(storage_get_files, note)

    db.update(note)

    return note


def storage_cleanup(uuid):
    storage.cleanup(uuid)


@app.delete("/api/note/{uuid}", response_model=Note)
async def api_note_delete(uuid: str, background_tasks: BackgroundTasks):
    background_tasks.add_task(storage_cleanup, uuid)
    return db.delete(uuid)


@app.put("/api/note/", response_model=Note)
async def api_note_new(item: Note, background_tasks: BackgroundTasks):

    note = item.dict()

    dt = datetime.today()
    seconds = int(dt.timestamp())

    note["created"] = seconds
    note["updated"] = seconds

    note = db.create(note)

    if storage_provider != "none":
        # note first has to be created, in  order to have it's ID/_id
        # and afterwards images will have to be parsed and downloaded
        # and note itself - updated.
        background_tasks.add_task(storage_get_files, note)

    return note


@app.get("/api/tags/")
async def api_tags(response: Response, query: str = ''):

    if db_conn_err:
        db_conn_err_persisting = connection()
        if db_conn_err_persisting:
            response.status_code = 503
            return {"error": db_conn_err_persisting}

    # list(set()) removes duplicates from list
    tags = list(set(db.tags()))

    tags = list(filter(lambda tag: tag != None, tags))

    # filtering tags
    if len(query) > 0:
        tags = list(filter(lambda tag: query.lower() in tag.lower(), tags))

    tags.sort()

    return tags


@app.get("/health")
async def health(request: Request):
    return {
        "status": "OK"
    }


@app.get("/{path:path}")
async def index(request: Request, path: str = ''):
    return templates.TemplateResponse("index.html", {"request": request, "app_version": package_json["version"]})
