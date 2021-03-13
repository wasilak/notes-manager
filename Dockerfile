# building stage
FROM python:3-slim as builder
RUN apt-get update \
  && apt-get install build-essential curl gnupg -y \
  && apt-get clean
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update \
  && apt-get install yarn -y \
  && apt-get clean
WORKDIR /app
COPY ./app /app
ENV PATH="/root/.local/bin:${PATH}"
RUN pip install --user -r requirements.txt
RUN yarn install

# production stage
FROM python:3-slim as app
COPY --from=builder /root/.local /root/.local
COPY --from=builder /app/ /app/
RUN apt-get update \
  && apt-get install curl -y \
  && apt-get clean
ENV PATH="/root/.local/bin:${PATH}"
WORKDIR /app
EXPOSE 5000
HEALTHCHECK --interval=5s --timeout=1s CMD curl -f http://localhost:5000/health || exit 1
CMD ["uvicorn", "main:app", "--host=0.0.0.0", "--port=5000", "--log-level=info"]
