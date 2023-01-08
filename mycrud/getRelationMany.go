package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type GetRelationManyImpl interface {
	GetRelationManyHandler()
	GetRelationManyDecode() kithttp.DecodeRequestFunc
	GetRelationManyEndpoint() endpoint.Endpoint
	GetRelationMany(ctx context.Context, id string, data interface{}) (err error)
}

type GetRelationMany struct {
	Repo *Repo
	*KitHttpConfig
}

type GetRelationManyRequest struct {
	Id string `json:"id"`
}

func (g *GetRelationMany) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetRelationManyRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		return req, nil
	}
}

func (g *GetRelationMany) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRelationManyRequest)
		res, err := g.Repo.GetRelationMany(ctx, req.Id, nil)
		return res, err
	}
}
