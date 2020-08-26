FROM golang:alpine

RUN mkdir /app

WORKDIR /app

COPY . /app/

RUN go build -o authservice .

CMD ["./authservice"]