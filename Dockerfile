FROM alpine

RUN apk add --no-cache go

WORKDIR /app

COPY ./go.sum ./go.mod ./

RUN go mod download

COPY *.go ./

ENV CGO_ENABLED=0

RUN go build -o ./otel_tracing_demoapp .

RUN chmod 500 ./otel_tracing_demoapp
RUN chown 1000:1000 ./otel_tracing_demoapp

FROM alpine

WORKDIR /app

USER 1000:1000

COPY --chown=1000:1000 --from=0 /app/otel_tracing_demoapp ./

ENTRYPOINT [ "/app/otel_tracing_demoapp" ]
