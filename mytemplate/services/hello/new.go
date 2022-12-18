package hello

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func NewMws(ms []endpoint.Middleware) (res Mws) {
	res = Mws{}
	for _, v := range ms {
		MethodAddMws(res, v, MethodNameList)
	}
	return
}

func NewService(base BaseService, log *zap.SugaredLogger) Service {
	var ms []Middleware
	ms = append(ms, NewLogging(log), NewTracing())

	for _, m := range ms {
		base = m(base)
	}
	return base
}

func NewOps(ops []http.ServerOption) (res Ops) {
	res = Ops{}
	for _, v := range ops {
		MethodAddOps(res, v, MethodNameList)
	}
	return
}

var HelloSet = wire.NewSet(
	MakeHTTPHandler,
	NewService,
	New,
	NewOps,
	NewMws,
)
