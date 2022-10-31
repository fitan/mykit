package myhttp

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type ErrorHandler struct {
	log *zap.SugaredLogger
}

func NewErrorHandler(log *zap.SugaredLogger) *ErrorHandler {
	return &ErrorHandler{log: log}
}

func (e *ErrorHandler) Handle(ctx context.Context, err error) {
	e.log.Errorw("error handler", "traceId", trace.SpanFromContext(ctx).SpanContext().TraceID().String(), "err", err.Error())
}

func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	ResponseJsonEncode(w, map[string]interface{}{"code": http.StatusInternalServerError, "err": err.Error(), "traceId": trace.SpanFromContext(ctx).SpanContext().TraceID().String()})
}
