package buffer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

func TestAddSpan(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	rb := NewRingBuffer(2*time.Second, 5, logger)

	span1 := ptrace.NewSpan()
	span1.SetTraceID([16]byte{1})
	span1.SetName("span1")
	span2 := ptrace.NewSpan()
	span2.SetTraceID([16]byte{2})
	span2.SetName("span2")

	// Add spans to the buffer
	err := rb.AddSpan(span1)
	assert.NoError(t, err)

	err = rb.AddSpan(span2)
	assert.NoError(t, err)

	spans, err := rb.GetSpans(span1.TraceID().String())
	assert.NoError(t, err)
	assert.Len(t, spans, 1) // Only one span should be associated with trace ID 1
	assert.Equal(t, "span1", spans[0].Name())
}

func TestAddSpanCapacity(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	rb := NewRingBuffer(1*time.Second, 2, logger)

	// Create test spans
	span1 := ptrace.NewSpan()
	span1.SetTraceID([16]byte{1})
	span1.SetName("span1")

	span2 := ptrace.NewSpan()
	span2.SetTraceID([16]byte{2})
	span2.SetName("span2")

	span3 := ptrace.NewSpan()
	span3.SetTraceID([16]byte{3})
	span3.SetName("span3")

	// Add spans to the buffer
	err := rb.AddSpan(span1)
	assert.NoError(t, err)

	err = rb.AddSpan(span2)
	assert.NoError(t, err)

	err = rb.AddSpan(span3)
	assert.NoError(t, err)

	// Now the buffer has 3 spans, but the capacity is 2, so the oldest span (span1) should be removed.

	// Retrieve spans for trace ID 1 (should be removed)
	spans, err := rb.GetSpans(span1.TraceID().String())
	assert.Error(t, err)
	assert.Nil(t, spans)

	// Retrieve spans for trace ID 2
	spans, err = rb.GetSpans(span2.TraceID().String())
	assert.NoError(t, err)
	assert.Len(t, spans, 1)
	assert.Equal(t, "span2", spans[0].Name())

	// Retrieve spans for trace ID 3
	spans, err = rb.GetSpans(span3.TraceID().String())
	assert.NoError(t, err)
	assert.Len(t, spans, 1)
	assert.Equal(t, "span3", spans[0].Name())
}
