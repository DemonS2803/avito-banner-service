FROM golang:1.22

WORKDIR /app
ADD go.mod .
COPY . .
RUN go build -o main main.go

ENTRYPOINT ["./main"]