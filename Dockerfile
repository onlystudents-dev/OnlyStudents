FROM oven/bun:alpine AS frontend
WORKDIR /app/frontend
COPY frontend/ ./
RUN bun install && bun run vite build

FROM golang:1.26.5-alpine AS serve
WORKDIR /go/src/app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o /go/bin/app .

FROM gcr.io/distroless/static-debian13
WORKDIR /app
COPY --from=frontend /app/frontend/dist /app/frontend/dist
COPY ./migrations /app/migrations
COPY --from=serve /go/bin/app /app/app
CMD ["/app/app"]
