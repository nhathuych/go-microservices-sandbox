FROM alpine:latest

WORKDIR /app

COPY mailerApp /app

CMD ["/app/mailerApp"]
