# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR $GOPATH/src/github.com/digitalxtian/golang/10_REST_api
COPY . . 
RUN go mod init
RUN go build -o main 
EXPOSE 8080
CMD ["./main"]

# run "docker build . --tag {container-name}"
# access conatiner by running "docker run --publish 8080:8080 api-server"
# curl hhtp://localhost:8080/