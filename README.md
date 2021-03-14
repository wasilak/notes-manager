# Notes Manager 
[![Total alerts](https://img.shields.io/lgtm/alerts/g/wasilak/notes-manager.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/wasilak/notes-manager/alerts/) [![Language grade: JavaScript](https://img.shields.io/lgtm/grade/javascript/g/wasilak/notes-manager.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/wasilak/notes-manager/context:javascript) [![Language grade: Python](https://img.shields.io/lgtm/grade/python/g/wasilak/notes-manager.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/wasilak/notes-manager/context:python) [![Maintainability](https://api.codeclimate.com/v1/badges/12f39774bcfc138889cb/maintainability)](https://codeclimate.com/github/wasilak/notes-manager/maintainability)

- markdown (github flavored)
- Python FastAPI as a backend
- quick start with: `docker-compose up -d --remove-orphans --build`
- uses `.env` for configuration:

  ```shell
  ELASTICSEARCH=localhost:9200
  DB_PROVIDER=elasticsearch
  # DB_PROVIDER=file
  ```

- database providers:
 - Elasticsearch (fully supportd, run via docker-compose or externally)
 - TinyDb, not fully supportd, for testing

- install dependencies with:
 - `pip install -U -r requirements.txt`
 - `yarn install`

- run locally with: ` uvicorn main:app --reload --log-level debug`
- create `notes` index with: `curl -s -XPUT http://localhost:9200/notes`
- [Growl notifications](http://jvandemo.github.io/angular-growl-notifications/)
