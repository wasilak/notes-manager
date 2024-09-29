FROM quay.io/wasilak/golang:1.23-alpine as builder
ARG VERSION=main

COPY . /app
WORKDIR /app/
RUN mkdir -p ./dist
RUN go build -ldflags "-X github.com/wasilak/notes-manager/libs/common.Version=${VERSION}" -o ./dist/notes-manager

FROM quay.io/wasilak/alpine:3

COPY --from=builder /app/dist/notes-manager /notes-manager

CMD ["/notes-manager"]
