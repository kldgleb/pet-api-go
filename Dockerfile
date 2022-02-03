FROM golang:1.16-alpine

WORKDIR /app

COPY ./ ./

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go build -o test-api ./cmd/main.go

EXPOSE 8080

CMD [ "./test-api" ]