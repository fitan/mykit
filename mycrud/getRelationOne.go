package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type GetRelationOne struct {
	Repo *Repo
	*KitHttpConfig
	RelationTableName string
}

func NewGetRelationOne(repo *Repo, kitHttpConfig *KitHttpConfig, relationTableName string) *GetRelationOne {
	return &GetRelationOne{Repo: repo, KitHttpConfig: kitHttpConfig, RelationTableName: relationTableName}
}

type GetRelationOneRequest struct {
	Id string `json:"id"`
}

func (g *GetRelationOne) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetRelationOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		return req, nil
	}
}

func (g *GetRelationOne) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRelationOneRequest)
		res, err := g.Repo.GetRelationOne(ctx, req.Id, g.RelationTableName, nil)
		return res, err
	}
}
