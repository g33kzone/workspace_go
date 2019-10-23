 FROM  golang:alpine3.9 as builder
 RUN apk update && apk add git

 WORKDIR /app

 COPY go.mod .
 COPY go.sum .

 RUN go mod download

 COPY . .

 RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

 FROM scratch
 COPY --from=builder /app/go_user /app/
 EXPOSE 8080

 ENTRYPOINT [ "/app/go_user" ]