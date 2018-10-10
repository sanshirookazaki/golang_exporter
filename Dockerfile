FROM golang:1.11.0
LABEL maintainer "s.okazaki"

WORKDIR /go/src/golang_exporter
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build

ENTRYPOINT ["./golang_exporter"]
