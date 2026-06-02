# Echo Service

A simple Go-based Echo and WebSocket service built with Gin.

## Setup

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run the server**:
   ```bash
   go run cmd/main.go
   ```

## Endpoints

- **HTTP Echo**: `ANY /echo`
  - Returns any payload (JSON, text, or query parameters) inside the `data` response field.
- **WebSocket Echo & Auto Ping-Pong**: `GET /websocket`
  - Automatically replies `"pong"` to `"ping"` messages (supports both text and binary).
  - Echoes other text or binary payloads as-is.

## Testing

Run unit tests:
```bash
go test -v ./internal/web/restapi/...
```
