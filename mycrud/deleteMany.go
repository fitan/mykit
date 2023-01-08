package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"strings"
)

type DeleteManyRequest struct {
	Ids []string `json:"ids"`
}

type DeleteMany struct {
	Repo *Repo
	*KitHttpConfig
}

func (d *DeleteMany) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := DeleteManyRequest{}
		ids := r.URL.Query().Get("ids")
		req.Ids = strings.Split(ids, ",")
		return req, nil
	}
}

func (d *DeleteMany) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteManyRequest)
		err = d.Repo.DeleteMany(ctx, req.Ids)
		return nil, err
	}
}
