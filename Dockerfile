FROM python:3-alpine

RUN apk add --update --no-cache build-base dumb-init yarn

COPY ./requirements.txt /requirements.txt

RUN pip install -r requirements.txt

COPY . /app/

WORKDIR /app

RUN yarn install

EXPOSE 5000

ENTRYPOINT ["/usr/bin/dumb-init", "--", "uvicorn", "main:app", "--host=0.0.0.0", "--port=5000"]

CMD ["--log-level=info"]
