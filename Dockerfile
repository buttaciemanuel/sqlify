FROM golang:1.23.6 AS builder

WORKDIR /app

COPY ./golang/src/ ./

RUN curl https://install.duckdb.org | sh
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o sqlify .

CMD ["./sqlify", "serve", "--configuration", "/data/config.yaml", "--port", "3001"]
