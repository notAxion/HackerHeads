FROM golang:1.15

WORKDIR '/app'

COPY go.mod go.sum ./
RUN go get -u ./...

#COPY src_github/. /usr/local/go/src/
COPY . .

RUN go build main.go
CMD ["./main"]