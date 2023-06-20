FROM golang:1.20

WORKDIR /app/go

COPY . /app/go

RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT ["go", "run", "main.go"]
