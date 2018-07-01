FROM python:3-alpine

COPY . /app

RUN apk add --update --no-cache nodejs nodejs-npm

WORKDIR /app

RUN npm install

RUN pip install -r requirements.txt

CMD ["flask", "run", "--host=0.0.0.0" ,"--with-threads", "--eager-loading"]
