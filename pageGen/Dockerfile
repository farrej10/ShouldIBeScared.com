FROM golang:1.17-alpine as builder

RUN apk update
RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /go/src/github.com/farrej10/ShouldIBeScared.com/

ADD pageGen/ ./
COPY movie/* movie/
COPY ../go.mod ./
COPY ../go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/pageGen

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/pageGen /go/bin/pageGen
COPY --from=builder /go/src/github.com/farrej10/ShouldIBeScared.com/templates/* /var/www/


EXPOSE 8085

CMD [ "/go/bin/pageGen" ]