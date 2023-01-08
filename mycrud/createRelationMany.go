package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

type CreateRelationMany struct {
	Repo *Repo
	*KitHttpConfig
}

type CreatRelationManyRequest struct {
	Id   string      `json:"id"`
	Body interface{} `json:"body"`
}

func (c *CreateRelationMany) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreatRelationManyRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		req.Body = c.Repo.RelationTableMsg.manyObjFn()
		err = json.NewDecoder(r.Body).Decode(&req.Body)
		if err != nil {
			err = errors.Wrap(err, "json.Decode")
			return
		}
		return req, nil
	}
}

func (c *CreateRelationMany) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreatRelationManyRequest)
		err = c.Repo.CreateRelationMany(ctx, req.Id, req.Body)
		return nil, err
	}
}
