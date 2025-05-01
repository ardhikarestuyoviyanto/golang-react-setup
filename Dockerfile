# ===== Stage 1: Build React Frontend =====
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Copy and install frontend dependencies
COPY app/views/js/package*.json ./
RUN npm install

# Copy all frontend source
COPY app/views/js ./

# Build React/Vite output
RUN npm run build

# ===== Stage 2: Build Go Backend =====
FROM golang:1.23-alpine AS backend 

WORKDIR /app

# Copy env file (pertimbangkan untuk menggunakan variabel environment di Railway)
COPY .env .env

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy entire backend source
COPY . ./

RUN go mod tidy

# Copy built frontend files into original JS dist path
COPY --from=frontend-builder /app/dist ./app/views/js/dist

# Build Go binary
RUN go build -o main .

# Expose port 1323
EXPOSE 1323

# Command to run the Go app
CMD ["./main"]
