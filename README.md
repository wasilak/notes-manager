# notes manager

- markdown (github flavored)
- Python Flask as a backend
- quick start with: `docker-compose up -d --remove-orphans --build`
- uses `.env` for configuration:
```
FLASK_APP=app.py
FLASK_RUN_PORT=5000
FLASK_ENV=development
ELASTICSEARCH=localhost:9200
DB_PROVIDER=elasticsearch
# DB_PROVIDER=file
```
- database providers:
 - Elasticsearch (fully supportd, run via docker-compose or externally)
 - TinyDb, not fully supportd, for testing
- install dependencies with:
 - `pip install -U -r requirements.txt`
 - `npm install`
- run locally with: `flask run -with-threads --eager-loading`
- create `notes` index with: `curl -s -XPUT http://localhost:9200/notes`
- [Growl notifications](http://jvandemo.github.io/angular-growl-notifications/)
