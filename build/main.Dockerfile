FROM golang:1.24.0-alpine AS builder

COPY . /github.com/totorialman/kr-net-6-sem
WORKDIR /github.com/totorialman/kr-net-6-sem

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main.go

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/totorialman/kr-net-6-sem/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8010

ENTRYPOINT ["./.bin"]