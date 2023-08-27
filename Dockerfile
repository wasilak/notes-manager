ARG VERSION=main
FROM quay.io/wasilak/golang:1.21-alpine as builder

ADD . /app
WORKDIR /app/
RUN mkdir -p ./dist
RUN go build -ldflags "-X github.com/wasilak/notes-manager/libs/common.Version=${VERSION}" -o ./dist/notes-manager

FROM quay.io/wasilak/alpine:3

COPY --from=builder /app/dist/notes-manager /notes-manager

CMD ["/notes-manager"]
