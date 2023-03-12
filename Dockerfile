FROM golang:1.18-alpine
RUN mkdir api 
WORKDIR /api
COPY ./ ./
RUN go mod tidy -compat=1.17
RUN go mod vendor
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]