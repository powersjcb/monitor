package tracer

import (
	"github.com/honeycombio/opentelemetry-exporter-go/honeycomb"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

const sampleFraction = 0.1

func InitTracer(hcAPIKey string) {
	exporter, err := honeycomb.NewExporter(
		honeycomb.Config{APIKey: hcAPIKey},
		honeycomb.TargetingDataset("monitor.jacobpowers.me"),
		honeycomb.WithServiceName("monitor.jacobpowers.me"),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.ProbabilitySampler(sampleFraction)}),
		sdktrace.WithSyncer(exporter),
	)
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}
