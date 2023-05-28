FROM golang:1.20

WORKDIR /worker
COPY . .

RUN go mod download && go mod verify

RUN go build -o main cmd/main.go

CMD ["/worker/main"]