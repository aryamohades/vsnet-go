FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOGC=100 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main sock/main.go

FROM scratch
COPY --from=builder /build/main /
CMD ["./main"]