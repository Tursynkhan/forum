FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o main ./cmd/main.go

FROM ubuntu:20.04
WORKDIR /app

LABEL "author"="https://01.alem.school/git/tursynkhan"
LABEL build_date="2023-01-28"

COPY --from=builder /app .
EXPOSE 8080
CMD ["./main"]