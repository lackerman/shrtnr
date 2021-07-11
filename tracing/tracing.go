package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var tracer = otel.Tracer("gin-server")

var httpClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}

// HttpRequest performs the desired request with tracing context
func HttpRequest(ctx context.Context, reqType, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/spantest", body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return httpClient.Do(req)
}

// InitTracer sets up an OpenTelemetry trace context with a Zipkin exporter
// The resources used to set this up can be found here:
// - https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
// - https://github.com/open-telemetry/opentelemetry-go/blob/main/example/zipkin/main.go
func InitTracer(url string) (func(), error) {
	exporter, err := zipkin.New(
		url,
		zipkin.WithSDKOptions(sdktrace.WithSampler(sdktrace.NeverSample())),
	)
	if err != nil {
		return nil, err
	}
	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("shrtnr"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return func() {
		_ = tp.Shutdown(context.Background())
	}, nil
}

func Middleware() gin.HandlerFunc {
	return otelgin.Middleware("shrtnr")
}

func NewSpan(log func(msg string, keysAndValues ...interface{}), c context.Context, handlerName string) (context.Context, trace.Span) {
	shortName := strings.TrimPrefix(handlerName, "github.com/lackerman/shrtnr/handlers.(*handler).")
	shortName = strings.Split(shortName, "-")[0]
	ctx, span := tracer.Start(c, handlerName)
	log(fmt.Sprintf("executing handler: %v", shortName),
		"trace_id", span.SpanContext().TraceID().String(),
		"span_id", span.SpanContext().SpanID().String(),
	)
	return ctx, span
}
