FROM golang
ENV GO111MODULE=on
ADD ./ /go/src/app
WORKDIR /go/src/app
# WORKDIR /app
#COPY ./go.mod .
#COPY ./go.sum .
#RUN go mod download
COPY . .
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
ENV PORT=3001
EXPOSE 3001
#ENTRYPOINT ["/app/httpserver"]
ENTRYPOINT ["go","run","main.go"]
#CMD ["go","run","main.go"]