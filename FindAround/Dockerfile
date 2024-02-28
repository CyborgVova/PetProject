FROM golang:1.21-alpine
ENV GOPATH="/"
WORKDIR $PWD/restaurant
COPY . .
EXPOSE 8888
RUN go build -o app main.go
CMD ["./app"]
