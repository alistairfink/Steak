FROM golang:1.12.6 as builder
LABEL maintainer="AlistairFink <alistairfink@gmail.com>"

WORKDIR /go/src/github.com/alistairfink/Steak/Frontend
COPY ./Frontend .
RUN GOOS=js GOARCH=wasm go build -o main.wasm -ldflags="-s -w"

FROM docker.io/library/nginx:latest

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /go/src/github.com/alistairfink/Steak/Frontend /usr/share/nginx/html
COPY ./bin .

CMD nginx && /Steak