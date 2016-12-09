FROM python:3-alpine

RUN apk --no-cache --update add git postgresql-dev gcc g++ make libffi-dev openssl-dev

COPY ./app /app/snipt

VOLUME /app/snipt

WORKDIR /app/snipt

RUN pip install -r requirements.txt

RUN python manage.py syncdb
