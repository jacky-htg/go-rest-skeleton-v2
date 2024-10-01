FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o go-server

FROM alpine:latest 
RUN apk add ca-certificates
COPY --from=builder /app/go-server /app/go-server
EXPOSE 8080
CMD ["/app/go-server"]