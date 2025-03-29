FROM golang:1.23-bookworm

WORKDIR /app

COPY . .

RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]