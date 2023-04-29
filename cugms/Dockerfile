FROM golang:1.20-alpine

RUN mkdir -p /app
WORKDIR /app

COPY . /app

RUN go mod download

RUN go build ./cmd/service/main.go

CMD [ "./main" ]
