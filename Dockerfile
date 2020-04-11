FROM golang:1.14.1

RUN mkdir -p /go/src/github.com/dinopuguh/kawulo-crawler/

WORKDIR /go/src/github.com/dinopuguh/kawulo-crawler/

COPY . .

RUN go build -o crawler main.go

EXPOSE 9090

CMD /go/src/github.com/dinopuguh/kawulo-crawler/crawler