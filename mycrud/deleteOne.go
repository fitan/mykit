package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type DeleteOne struct {
	Repo *Repo
	*KitHttpConfig
}

type DeleteOneRequest struct {
	Id string `json:"id"`
}

func (d *DeleteOne) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := DeleteOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		return req, nil
	}
}

func (d *DeleteOne) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteOneRequest)
		err = d.Repo.DeleteOne(ctx, req.Id)
		return nil, err
	}
}
