#FROM golang:alpine
#
#WORKDIR /
#ENV PORT=8000
#CMD["go","run","main.go"]

#FROM golang:alpine
#ADD ./src /go/src/app
#WORKDIR /go/src/app
#ENV PORT=3001
##CMD["go","build", "./src/main.go"]
#CMD ["go", "run", "main.go"]

FROM golang
ENV GO111MODULE=on
ADD ./ /app
WORKDIR /app
COPY . .
RUN go mod download
RUN go build webhook.go
EXPOSE 8082
ENTRYPOINT ["go","run","webhook.go"]
