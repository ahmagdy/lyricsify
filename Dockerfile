FROM golang:1.13.3-alpine  as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lyricsify .

# final stage
FROM scratch
COPY --from=builder /app/lyricsify /app/
ENTRYPOINT ["/app/lyricsify"]
