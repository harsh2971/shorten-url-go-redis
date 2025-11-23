FROM golang:alpine as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go mod download
RUN go build -o main .

#stage 2

FROM alpine:latest

RUN mkdir -p /app && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser -h /app appuser

COPY --from=builder /build/main /app/

RUN chown -R appuser:appuser /app

USER appuser

WORKDIR /app

EXPOSE 3000

CMD ["./main"]