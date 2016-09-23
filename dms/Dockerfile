FROM alpine:latest
RUN apk add --update sqlite
RUN mkdir /app
WORKDIR /app
ADD app /app
EXPOSE 8080
CMD ["/app"]

