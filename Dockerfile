FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY *.pem ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tv-bot

CMD [ "/tv-bot" ]