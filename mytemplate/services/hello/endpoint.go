package hello

import (
	"context"

	myhttp "github.com/fitan/mykit/myhttp"
	endpoint "github.com/go-kit/kit/endpoint"
)

const HelloMethodName = "Hello"

var MethodNameList = []string{HelloMethodName}

type Endpoints struct {
	HelloEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{HelloEndpoint: makeHelloEndpoint(s)}
	for _, m := range dmw[HelloMethodName] {
		eps.HelloEndpoint = m(eps.HelloEndpoint)
	}

	return eps
}
func makeHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HelloRequest)
		res, err := s.Hello(ctx, req.ID, req.Query)
		return myhttp.WrapResponse(res, err)

	}
}

type Mws map[string][]endpoint.Middleware

func MethodAddMws(mw Mws, m endpoint.Middleware, methods []string) {
	for _, v := range methods {
		mw[v] = append(mw[v], m)
	}
}
