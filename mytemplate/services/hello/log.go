package hello

import (
	"context"
	"time"

	trace "go.opentelemetry.io/otel/trace"
	zap "go.uber.org/zap"
)

type logging struct {
	logger *zap.SugaredLogger
	next   Service
}

func (s *logging) Hello(ctx context.Context, id string, query Query) (res HelloResponse, err error) {
	defer func(begin time.Time) {
		if err != nil {
			s.logger.Errorw("Hello error", "error", err, "id", id, "query", query, "took", time.Since(begin), "traceId", trace.SpanContextFromContext(ctx).TraceID().String())
		} else {
			s.logger.Infow("Hello success", "id", id, "query", query, "took", time.Since(begin), "traceId", trace.SpanContextFromContext(ctx).TraceID().String())
		}
	}(time.Now())
	return s.next.Hello(ctx, id, query)
}
func NewLogging(logger *zap.SugaredLogger) Middleware {
	return func(next Service) Service {
		return &logging{logger: logger, next: next}
	}
}
