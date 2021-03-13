# Notes Manager

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
