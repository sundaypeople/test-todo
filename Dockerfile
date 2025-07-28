FROM golang:1.24 AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -o todo main.go

FROM scratch AS prod
COPY --from=builder /app/todo /todo
ENTRYPOINT ["/todo"]
EXPOSE 8080
