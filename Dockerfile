FROM golang:alpine
RUN apk add --update git
WORKDIR /go/src/github.com/hagbarth/GoChat
ADD . /go/src/github.com/hagbarth/GoChat
RUN go get github.com/satori/go.uuid github.com/gorilla/websocket
RUN go build main.go
CMD ["./main"]
