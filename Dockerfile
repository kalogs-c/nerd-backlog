FROM cgr.dev/chainguard/go AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /out/nerd-backlog ./cmd/http

FROM cgr.dev/chainguard/static

WORKDIR /app

COPY --from=builder /out/nerd-backlog /app/nerd-backlog

ENV HTTP_HOST=0.0.0.0 \
    HTTP_PORT=42069

EXPOSE 42069

ENTRYPOINT ["/app/nerd-backlog"]
