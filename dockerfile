FROM golang:latest

WORKDIR /server

COPY ./entrypoint.sh /entrypoint.sh

ADD go.mod go.mod

RUN go mod download

ADD . .

EXPOSE "${API_INT_PORT}"

ENTRYPOINT ["sh", "/entrypoint.sh"]
