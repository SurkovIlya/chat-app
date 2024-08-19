FROM golang:1.22.2-alpine3.19 as builder

COPY go.mod go.sum /go/src/github.com/SurkovIlya/chat-app/
WORKDIR /go/src/github.com/SurkovIlya/chat-app
RUN go mod download
COPY . /go/src/github.com/SurkovIlya/chat-app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/chat-app github.com/SurkovIlya/chat-app


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/SurkovIlya/chat-app/build/chat-app /usr/bin/chat-app

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/chat-app"]