FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch
ADD www /www
ADD bin/atlas-mapserver /atlas-mapserver
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/atlas-mapserver"]