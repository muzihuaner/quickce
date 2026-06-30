# ---- Stage 1: Build frontend ----
FROM node:22-alpine AS frontend
WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# ---- Stage 2: Build Go backend ----
FROM golang:1.26-alpine AS backend
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o quickce-server ./cmd/server

# ---- Stage 3: Runtime ----
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend /build/quickce-server .
COPY --from=frontend /frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./quickce-server"]
