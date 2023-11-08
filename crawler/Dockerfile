# ===== build stage ====
FROM golang:1.20.10-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-w -s" -o main ./crawler/main.go

# ===== deploy stage ====
FROM golang:1.20.10-alpine

RUN apt update

RUN apk add chromium

COPY --from=builder /app/main /app/main

CMD ["/app/main"]
