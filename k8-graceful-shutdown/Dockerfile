FROM golang:1.22-bullseye as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app

EXPOSE 8888

CMD ["./app"]
