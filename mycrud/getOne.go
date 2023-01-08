package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type GetOneImpl interface {
	GetOneHandler()
	GetOneDecode() kithttp.DecodeRequestFunc
	GetOneEndpoint() endpoint.Endpoint
	GetOne(ctx context.Context, tableName, id string) (data interface{}, err error)
}

type GetOne struct {
	Repo *Repo
	*KitHttpConfig
}

type GetOneRequest struct {
	Id string `json:"id"`
}

func (g *GetOne) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		return req, nil
	}
}

func (g *GetOne) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOneRequest)
		res, err := g.Repo.GetOne(ctx, req.Id)
		return res, err
	}
}
