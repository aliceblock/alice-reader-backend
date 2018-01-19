FROM golang:1.9.2-alpine AS go-build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

##############################

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=go-build /go/src/app .
CMD ["./app"]