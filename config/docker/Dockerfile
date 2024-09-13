FROM golang:latest

LABEL authors="Artem Gafarov"
LABEL email="gafarov@realnoevremya.ru"
LABEL version="2.0"

WORKDIR /src/text2text-service
ADD . .

RUN go mod download && go mod verify

EXPOSE 80

WORKDIR /src/text2text-service/cmd/text2text-service
RUN go build

CMD ["./text2text-service"]