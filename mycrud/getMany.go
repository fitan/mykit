package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

type GetMany struct {
	Repo *Repo
	*KitHttpConfig
}

type GetManyData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type GetManyRequest struct {
}

func (g *GetMany) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		return nil, nil
	}
}

func (g *GetMany) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := g.Repo.GetMany(ctx, nil)
		return res, err
	}
}
