FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./ ./

RUN go build -o bin/api cmd/api/main.go

EXPOSE 3000

CMD [ "bin/api" ]