package trace

import (
	"context"

	"github.com/stackify/stackify-go-apm/config"
	"github.com/stackify/stackify-go-apm/trace/span"
	"github.com/stackify/stackify-go-apm/transport"
	"github.com/stackify/stackify-go-apm/utils"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type StackifySpanExporter struct {
	c *config.Config
	t transport.Transport
}

// ExportSpans method convert spans into stackify span format.
func (ssp *StackifySpanExporter) ExportSpans(ctx context.Context, spans []*export.SpanData) error {
	trace, err := ssp.toStackifyTrace(spans)
	if err != nil {
		return nil
	}
	ssp.t.HandleTrace(trace)
	return nil
}

// Shutdown method ensures we send all valid spans from queue.
func (ssp *StackifySpanExporter) Shutdown(context.Context) error {
	ssp.t.SendAll()
	return nil
}

// toStackifyTrace method converts spans to stackify trace format.
func (ssp *StackifySpanExporter) toStackifyTrace(sd []*export.SpanData) (*span.StackifySpan, error) {
	stackifySpans := []*span.StackifySpan{}
	stackifySpansMap := make(map[string]*span.StackifySpan)
	stackifySpan := &span.StackifySpan{}

	// map spans by span id for easy access
	for _, s := range sd {
		stackifySpan := span.NewSpan(ssp.c, s)
		stackifySpans = append(stackifySpans, &stackifySpan)
		stackifySpansMap[utils.SpanIdToString(s.SpanContext.SpanID[:])] = &stackifySpan
	}

	// arrange spans by adding child span into parent stacks attribute
	for _, s := range stackifySpans {
		if s.ParentId != s.Id && s.ParentId != utils.SpanIdToString(span.InvalidSpanId[:]) {
			stackifySpansMap[s.ParentId].Stacks = append(stackifySpansMap[s.ParentId].Stacks, stackifySpansMap[s.Id])
		} else {
			stackifySpan = s
		}
	}
	return stackifySpan, nil
}

// NewStackifySpanExporter function creates a StackifySpanExporter.
func NewStackifySpanExporter(c *config.Config, t *transport.Transport) *StackifySpanExporter {
	return &StackifySpanExporter{
		c: c,
		t: *t,
	}
}
