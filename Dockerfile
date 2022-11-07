FROM golang:latest AS buildstage
RUN mkdir -p /go/src/NewsFeed
WORKDIR /go/src/NewsFeed
COPY ./ ./
RUN go env -w GO111MODULE=auto
RUN go build -o server ./cmd/server.go

FROM alpine:latest
WORKDIR /NewsFeed/cmd/
COPY --from=buildstage /go/src/NewsFeed .
RUN apk add libc6-compat
RUN mv server ./cmd/
CMD ["./server"]
EXPOSE 8080