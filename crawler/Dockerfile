# ===== build stage ====
FROM golang:1.20.10-bullseye as builder

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
FROM golang:1.20.10-alpine as deploy

WORKDIR /app

RUN apk update

RUN apk add chromium
RUN apk add libc6-compat
RUN ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2

COPY --from=builder /app/crawler/bin/main ./main

CMD ["/app/main"]
