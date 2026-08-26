FROM oven/bun:alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend/ .
RUN bun run vite build

FROM golang:1.26.5-alpine AS serve
WORKDIR /go/src/app
ENV CGO_ENABLED=0

COPY go.mod go.sum* ./
RUN go mod download

COPY . ./
RUN go build -o /go/bin/app .

FROM gcr.io/distroless/static-debian13
WORKDIR /app
COPY --from=frontend /app/frontend/dist /app/frontend/dist
COPY --from=serve /go/bin/app /app/app
CMD ["/app/app"]