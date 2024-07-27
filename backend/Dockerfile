FROM golang:1.22.5-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /daylist_backend

EXPOSE 8080

CMD ["/daylist_backend"]