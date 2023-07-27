FROM --platform=linux/amd64 golang:1.20.6-alpine3.18 AS builder

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE="on"

WORKDIR $GOPATH/src/github.com/diegosepusoto/basic-website-bff

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o $GOBIN/main src/main.go

FROM --platform=linux/amd64 alpine:3.18.2

WORKDIR /app

COPY --from=builder /go/bin .

EXPOSE 8082

CMD ["./main"]