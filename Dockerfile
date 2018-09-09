FROM golang:latest

RUN mkdir /go/src/app
ADD . /go/src/app

RUN go get -v -d ./...

RUN go build -v ./...

CMD ["substitutes-notify"]
