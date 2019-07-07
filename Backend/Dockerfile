FROM golang:1.12.6 as builder
LABEL maintainer="AlistairFink <alistairfink@gmail.com>"

WORKDIR /go/src/github.com/alistairfink/Steak/Backend
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/Steak .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/Steak/ .
COPY --from=builder /go/src/github.com/alistairfink/Steak/Backend/config.json .

EXPOSE 41691

CMD ["./Steak"] 