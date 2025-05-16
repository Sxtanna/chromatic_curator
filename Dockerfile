FROM golang:1.24-alpine3.21 AS builder
LABEL authors="Sxtanna"

ENV GOFLAGS="-mod=readonly"

RUN apk add --update --no-cache ca-certificates make git curl mercurial

RUN mkdir -p /workspace
WORKDIR /workspace
COPY . ./

ARG GOPROXY

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./curator ./cmd/...

FROM gcr.io/distroless/static-debian12

COPY --from=builder /workspace/curator /
COPY --from=builder /workspace/.env /

CMD ["/curator"]