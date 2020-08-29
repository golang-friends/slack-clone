FROM golang:alpine

# workdir performs mkdir & cd
WORKDIR /app

# These are cached so that for the code change,
# we don't redownload the dependencies unless go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o chat .
CMD ["./chat"]
