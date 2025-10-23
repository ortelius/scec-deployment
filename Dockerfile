FROM cgr.dev/chainguard/go@sha256:56cc6e94e74d1e30bf923bcfbcfdddb2fc4d4c7fb85f4fd0b7603c85dbf1c345 AS builder

WORKDIR /app
COPY . /app

RUN go mod tidy; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:dfe55497b1b74855e14b27f1710bf9658ebca69cbebe00e5370ab3bf6da2f9d1

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
