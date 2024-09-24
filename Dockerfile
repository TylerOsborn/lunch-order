FROM golang:1.22.1 AS backend-builder

WORKDIR /lunch-order

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux go build -o main

FROM node:20 AS frontend-builder

WORKDIR /lunch-order/frontend

COPY frontend/package*.json ./

RUN npm install

COPY frontend .

RUN npm run build

FROM ubuntu:22.04

WORKDIR /root/

RUN mkdir -p /root/database

COPY --from=backend-builder /lunch-order/main .

COPY --from=frontend-builder /lunch-order/frontend/dist ./frontend/dist

VOLUME /root/database

EXPOSE 8080:8080

CMD ["./main"]