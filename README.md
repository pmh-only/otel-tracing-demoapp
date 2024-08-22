# otel-tracing-demoapp

## Environment Variables
* `NODE_NAME`: name for this node (default: `application`)

* `NODE_LISTEN`: listening host:port for this node (default: `0.0.0.0:3000`)

* `NEXT_NODE`: host:port combinations for next nodes\
  First item is more weighted than next item.\
  (example: `service-a.default.svc:3000,service-b.default.svc:3000`)

* `OTEL_EXPORTER_OTLP_ENDPOINT`: otel endpoint for tracing
