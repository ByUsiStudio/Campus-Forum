# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM alpine:latest

RUN apk add --no-cache busybox-extras

WORKDIR /app

COPY --from=builder /app/dist /app

EXPOSE 80

CMD ["busybox", "httpd", "-f", "-p", "80", "-h", "/app"]