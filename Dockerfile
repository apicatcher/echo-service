FROM --platform=linux/amd64 docker.io/golang:1.23.10 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 go build -o server cmd/main.go

FROM --platform=linux/amd64 docker.io/alpine:3.18
WORKDIR /usr/local/app
COPY --from=builder /workspace/server .
COPY config.toml config.toml
RUN echo "./server" > start_up.sh
RUN chmod +x start_up.sh
CMD ["sh", "start_up.sh"]