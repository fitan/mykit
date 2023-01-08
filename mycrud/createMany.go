package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

type CreateMany struct {
	Repo *Repo
	*KitHttpConfig
}

type CreateManyRequest struct {
	TableName string `json:"tableName"`
	Body      interface{}
}

func (c *CreateMany) GetDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateManyRequest{}

		body := c.Repo.TableMsg.manyObjFn()
		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return
		}

		req.Body = body
		return req, nil
	}
}

func (c *CreateMany) GetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateManyRequest)
		err = c.Repo.CreateMany(ctx, req.Body)
		return nil, err
	}
}
