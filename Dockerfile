FROM golang:1.20.0-bullseye

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY cmd cmd/
COPY pkg pkg/
COPY Makefile .

RUN make build

CMD [ "/app/bin/rate-limit-server" ]