FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/meeting-svc ./cmd/meeting-events

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /out/meeting-svc /app/meeting-svc
COPY config /app/config
COPY internal/adapters/postgres/init /app/internal/adapters/postgres/init
ENV CONFIG_PATH=/app/config/prod.yaml
EXPOSE 8081
ENTRYPOINT ["/app/meeting-svc"]
