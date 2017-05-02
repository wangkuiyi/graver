FROM golang:wheezy

RUN apt-get update
RUN apt-get install -y graphviz
RUN apt-get clean -y

COPY . /go/src/github.com/wangkuiyi/graver
RUN go install github.com/wangkuiyi/graver

EXPOSE 8080
ENTRYPOINT ["graver"]
