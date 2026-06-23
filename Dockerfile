FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY src/ ./src/
COPY references.bib .
COPY assets/ ./assets/
RUN go build -C src -mod=vendor -o ../compiler
RUN ./compiler -path references.bib > index.html

FROM nginx:alpine
COPY config/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/assets /usr/share/nginx/html/assets
COPY --from=builder /app/index.html /usr/share/nginx/html/index.html
