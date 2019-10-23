 # Default to Go 1.11
 ARG GO_VERSION=1.12.0

 # First stage: build the executable.
 FROM golang:${GO_VERSION}-alpine AS builder

 #FROM golang:1.12

 RUN apk update && apk add git

 WORKDIR /app

 COPY go.mod .
 COPY go.sum .

 RUN go mod download

 COPY . .

 RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

 FROM scratch
 COPY --from=builder /app/go-bdd-ginkgo /app/
 EXPOSE 8080

 ENTRYPOINT [ "/app/go-bdd-ginkgo" ]