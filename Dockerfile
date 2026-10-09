# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY ui ./ui
RUN CGO_ENABLED=0 go build -o /out/scvis ./cmd/web

# Runtime stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/scvis ./scvis
COPY visual/scvis ./visual/scvis
ENV SERVER_PROTOCOL=http \
    LOG_DIR=stdout \
    PORT=8080
EXPOSE 8080
CMD ["./scvis"]
