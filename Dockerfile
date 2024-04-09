FROM golang:1.22

WORKDIR /app
ADD go.mod .
COPY . .
RUN go build -o main cmd/main.go

ENTRYPOINT ["./main"]