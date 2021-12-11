FROM python:3.9-alpine as builder

RUN apk --update --no-cache add yarn cargo build-base
WORKDIR /app
COPY ./app /app
ENV PATH="/root/.local/bin:${PATH}"
RUN pip install --user -r requirements.txt
RUN yarn install

# production stage
FROM python:3.9-alpine as app
COPY --from=builder /root/.local /root/.local
COPY --from=builder /app/ /app/

RUN apk --update --no-cache add curl

ENV PATH="/root/.local/bin:${PATH}"
WORKDIR /app
EXPOSE 5000
HEALTHCHECK --interval=5s --timeout=1s CMD curl -f http://localhost:5000/health || exit 1
CMD ["uvicorn", "main:app", "--host=0.0.0.0", "--port=5000", "--log-level=info"]
