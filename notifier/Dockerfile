FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN apk --no-cache add ca-certificates

RUN go build -o notifier .
CMD ["./notifier"]
