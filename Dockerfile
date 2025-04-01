FROM alpine:3.16 AS root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser -u 1001 -D -G app -h /home/app app

FROM golang:1.24 AS builder
WORKDIR /youtube-api-files
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY . . 
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./youtube-stats ./app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./youtube-stats ./app

FROM scratch AS final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder /youtube-api-files/youtube-stats /youtube-stats
USER app
ENTRYPOINT [ "/youtube-stats" ]
