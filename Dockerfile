FROM golang:latest

RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

RUN go get -v -d ./...

RUN go build -v ./...
RUN go install -v .

CMD ["app", "-v"]
