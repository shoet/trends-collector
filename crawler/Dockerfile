# ===== build stage ====
FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY ./crawler/go.mod ./crawler/go.mod
COPY ./crawler/go.sum ./crawler/go.sum
COPY ./go.mod ./go.mod

RUN --mount=type=cache,target=/go-mod-cache \
    cd ./crawler && go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    cd ./crawler && \
    go build -trimpath -ldflags="-w -s" -tags timetzdata -o ./bin/main ./cmd/crawltask/main.go

# ===== deploy stage ====
FROM golang:1.21-bullseye as deploy

RUN apt update -y
RUN apt install -y chromium

WORKDIR /app

COPY --from=builder /app/crawler/bin/main ./main

ENV CRAWLER_BROWSER_PATH=/usr/bin/chromium

CMD ["/app/main"]
