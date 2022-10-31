package hello

import (
	"context"
	"encoding/json"

	otel "go.opentelemetry.io/otel"
	attribute "go.opentelemetry.io/otel/attribute"
	codes "go.opentelemetry.io/otel/codes"
	trace "go.opentelemetry.io/otel/trace"
)

type tracing struct {
	next Service
}

func (s *tracing) Hello(ctx context.Context, id string, query Query) (res HelloResponse, err error) {
	_, span := otel.Tracer("services").Start(ctx, "Hello")

	defer func() {
		paramIn := map[string]interface{}{
			"id":    id,
			"query": query,
		}
		paramInJsonB, _ := json.Marshal(paramIn)
		span.AddEvent("paramIn", trace.WithAttributes(attribute.String("param list", string(paramInJsonB))))
		if err != nil {
			span.SetStatus(codes.Error, "Hello error")
			span.RecordError(err)
		}
		span.End()
	}()

	return s.next.Hello(ctx, id, query)
}
func NewTracing() Middleware {
	return func(next Service) Service {
		return &tracing{next: next}
	}
}
