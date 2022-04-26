FROM golang:1.18 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY internal/ ./internal/

# Create a staticly linked binary
ENV CGO_ENABLED=0
RUN go build -o bin/server cmd/server/main.go

FROM alpine as run

WORKDIR /app

COPY --from=build /app/bin/server /app/server

ENTRYPOINT [ "/app/server" ]