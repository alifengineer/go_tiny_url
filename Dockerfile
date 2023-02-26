FROM golang:1.18 as builder


ENV $GOPATH=/go
ENV $PATH=$GOPATH/bin:$PATH

#
RUN mkdir -p $GOPATH/src/github.dilmurodov/app
WORKDIR $GOPATH/src/github.dilmurodov/app

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
# RUN export CGO_ENABLED=0 && \
#     export GOOS=linux && \
#     go mod vendor && \
#     make build && \
#     mv ./bin/app /

COPY ./ ./

RUN go build -o /app ./cmd/main.go

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN migrate -path ./migrations/postgres \
    -database="postgres://postgres:admin@postgres-server:5454/shortener_db?sslmode=disable&x-migrations-table=migrations" up

FROM alpine
COPY --from=builder app .
ENTRYPOINT ["/app"]
