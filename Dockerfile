FROM quay.io/wasilak/python:3-slim

RUN apt-get update && apt-get install -y curl

WORKDIR /app
COPY ./app /app
RUN pip install --user -U -r requirements.txt

EXPOSE 5000
HEALTHCHECK --interval=5s --timeout=1s CMD curl -f http://localhost:5000/health || exit 1
CMD ["uvicorn", "main:app", "--host=0.0.0.0", "--port=5000", "--log-level=info"]
