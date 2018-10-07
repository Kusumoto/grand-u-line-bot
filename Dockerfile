FROM golang:1.11
WORKDIR /go/src/github.com/kusumoto/grand-u-line-bot
ADD ./ .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -v -vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/kusumoto/grand-u-line-bot .
CMD ["./main"]  