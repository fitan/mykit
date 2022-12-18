package hello

import (
	"context"
	"fmt"
)

type Middleware func(service Service) Service

// @tags hello
// @impl
type Service interface {
	// @kit-http /hello/{id} GET
	// @kit-http-request HelloRequest
	Hello(ctx context.Context, id string, query Query) (res HelloResponse, err error)
}

type service struct {
}

func (s *service) Hello(ctx context.Context, id string, query Query) (res HelloResponse, err error) {
	res.ID = id
	// @call query
	q := queryDTO(query)
	fmt.Println(q)
	return
}

type BaseService Service

func New() BaseService {
	return &service{}
}
