FROM golang:1.21-alpine
ENV GOPATH="/"
WORKDIR shortlink
COPY . .
EXPOSE 8080
RUN go build cmd/app/app.go
CMD ["./app", " -s inmemory"] 