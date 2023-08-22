# Notes Manager
[![Docker Repository on Quay](https://quay.io/repository/wasilak/notes-manager/status "Docker Repository on Quay")](https://quay.io/repository/wasilak/notes-manager) [![Maintainability](https://api.codeclimate.com/v1/badges/12f39774bcfc138889cb/maintainability)](https://codeclimate.com/github/wasilak/notes-manager/maintainability)

- markdown (github flavored)
- Go backend
- quick start with: `docker-compose up -d --remove-orphans --build`
- uses `.env` for configuration:

  ```shell
  DB_PROVIDER=mongodb
  ```

- database providers:
 - MongoDB 4.x

- install dependencies with:
 - `go mod tidy`
 - `yarn install`

- run locally with: ` uvicorn main:app --reload --log-level debug`
- create `notes` index with: `curl -s -XPUT http://localhost:9200/notes`
- [Growl notifications](http://jvandemo.github.io/angular-growl-notifications/)
