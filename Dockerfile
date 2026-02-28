# Stage 1: Frontend build
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Go build
FROM golang:1.24-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend/ ./backend/
COPY --from=frontend /app/frontend/dist ./backend/static/
RUN cd backend && CGO_ENABLED=1 go build -o /homemenu .

# Stage 3: Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=backend /homemenu /app/homemenu
WORKDIR /app
RUN mkdir -p /app/data/uploads
EXPOSE 8080
CMD ["./homemenu"]
