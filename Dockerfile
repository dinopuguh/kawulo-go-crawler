FROM golang:1.14.1

WORKDIR /app

EXPOSE 9090

COPY crawler .

CMD /app/crawler