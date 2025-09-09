FROM cgr.dev/chainguard/go@sha256:da99a150cc0b122a55e9dff7e696407d2e76027502129fb56f5dabbd767d1bd7 AS builder

WORKDIR /app
COPY . /app

RUN go mod tidy; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:722bbacce71a60ecc28ef72b64fa5d588591df558075a7d108c68be34bae03c7

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs docs

ENV ARANGO_HOST localhost
ENV ARANGO_USER root
ENV ARANGO_PASS rootpassword
ENV ARANGO_PORT 8529
ENV MS_PORT 8080

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
