FROM python:3-alpine

COPY . /app

ENV FLASK_ENV=production
ENV FLASK_RUN_PORT=5000
ENV FLASK_DEBUG=False
ENV FLASK_APP=app.py

RUN apk add --update --no-cache yarn

WORKDIR /app

RUN yarn install

RUN pip install -U pip

RUN pip install -r requirements.txt

CMD ["flask", "run", "--host=0.0.0.0" ,"--with-threads", "--eager-loading"]
