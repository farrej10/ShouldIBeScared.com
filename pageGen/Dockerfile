FROM golang:1.17-alpine

RUN apk update
RUN apk add git

WORKDIR /go/src/github.com/farrej10/ShouldIBeScared.com/

COPY ../go.mod ./
COPY ../go.sum ./
RUN go mod download
RUN go get -v ./...

COPY ./ ./



RUN go build -o ./pageGen/pageGen ./pageGen/pageGen.go 
EXPOSE 8085

CMD [ "./pageGen/pageGen" ]