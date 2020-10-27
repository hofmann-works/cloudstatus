FROM golang:1.15.3-alpine3.12 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o cloudstatus .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/cloudstatus /app/
WORKDIR /app
CMD ["./cloudstatus"]