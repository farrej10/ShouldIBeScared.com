FROM golang:1.17-alpine

RUN apk update
RUN apk add git

WORKDIR /go/src/github.com/farrej10/ShouldIBeScared.com/

COPY ../go.mod ./
COPY ../go.sum ./
RUN go mod download


COPY ./ ./


RUN go get -v ./...
RUN go build -o ./movieService/movieService ./movieService/movieService.go 
EXPOSE 50051

CMD [ "./movieService/movieService" ]