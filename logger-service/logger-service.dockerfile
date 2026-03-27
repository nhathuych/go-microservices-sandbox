FROM alpine:latest

WORKDIR /app

COPY loggerApp /app

CMD ["/app/loggerApp"]
