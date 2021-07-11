module github.com/lackerman/shrtnr

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-logr/logr v0.4.0
	github.com/syndtr/goleveldb v1.0.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.21.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/zipkin v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	k8s.io/klog/v2 v2.9.0
)
