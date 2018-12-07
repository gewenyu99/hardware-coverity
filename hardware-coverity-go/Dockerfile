FROM golang:latest

ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY . $GOPATH/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/
WORKDIR $GOPATH/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/
RUN dep ensure -vendor-only
RUN go build -o main .
CMD ["./main"]