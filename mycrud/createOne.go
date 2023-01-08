package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

type CreateOne struct {
	Repo *Repo
	*KitHttpConfig
}

type CreateOneRequest struct {
	Body interface{}
}

func (c *CreateOne) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateOneRequest{}

		body := c.Repo.TableMsg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return nil, err
		}

		req.Body = body
		return req, err
	}
}

func (c *CreateOne) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateOneRequest)
		err = c.Repo.CreateOne(ctx, req.Body)
		return nil, err
	}
}
