FROM python:3-alpine

COPY . /app

WORKDIR /app

RUN pip install -r requirements.txt

CMD ["flask", "run", "--host=0.0.0.0" ,"--with-threads", "--eager-loading"]
