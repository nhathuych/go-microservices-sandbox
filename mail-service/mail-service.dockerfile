FROM alpine:latest

WORKDIR /app

COPY mailerApp /app
COPY templates /app/templates

CMD ["/app/mailerApp"]
