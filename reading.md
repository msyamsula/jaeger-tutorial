-- trace propagation (different framework result in different propagation, see example in the package)
https://opentelemetry.io/docs/instrumentation/go/manual/#propagators-and-context
https://opentelemetry.io/docs/instrumentation/go/libraries/#available-packages
https://opentelemetry.io/ecosystem/registry/?language=go&component=instrumentation


-- jaeger alternative
https://zipkin.io/
https://github.com/openzipkin/zipkin/tree/master/docker/examples

folder "mytracer" in each service is kinda the same, need to refactor so it comply with DRY paradigm