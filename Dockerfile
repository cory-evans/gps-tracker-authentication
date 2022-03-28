FROM golang:1.18 as build

WORKDIR /app

RUN BIN="/usr/local/bin" && \
	VERSION="1.1.0" && \
	curl -sSL \
	"https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" \
	-o "${BIN}/buf" && \
	chmod +x "${BIN}/buf"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1 && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY buf.gen.yaml ./
COPY buf.lock ./
COPY buf.yaml ./
RUN buf mod update

COPY pkg/ ./pkg/

RUN buf generate

COPY internal/ ./internal/
COPY cmd/ ./cmd/

# Create a staticly linked binary
ENV CGO_ENABLED=0
RUN go build -o bin/server cmd/server/main.go

FROM alpine as run

WORKDIR /app

COPY --from=build /app/bin/server /app/server

EXPOSE 8080

ENTRYPOINT [ "/app/server" ]