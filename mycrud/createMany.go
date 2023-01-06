package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"net/http"
)

type CreateManyImpl interface {
	CreateManyHandler()
	CreateManyDecode() kithttp.DecodeRequestFunc
	CreateManyEndpoint() endpoint.Endpoint
	CreateMany(ctx context.Context, data interface{}) (err error)
}

type CreateMany struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (c *CreateMany) CreateManyHandler() {
	c.Crud.Handler(CreateManyMethodName, http.MethodPost, "/"+c.TableMsg.schema.Table+"/many", c.CreateManyEndpoint(), c.CreateManyDecode())
}

type CreateManyRequest struct {
	TableName string `json:"table_name"`
	Body      interface{}
}

func (c *CreateMany) CreateManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateManyRequest{}

		body := c.TableMsg.manyObjFn()
		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return
		}

		req.Body = body
		return req, nil
	}
}

func (c *CreateMany) CreateManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateManyRequest)
		err = c.CreateMany(ctx, req.Body)
		return nil, err
	}
}

func (c *CreateMany) CreateMany(ctx context.Context, data interface{}) (err error) {

	db, commit := c.Crud.db.Tx(ctx)
	defer commit(err)

	err = db.Model(c.TableMsg.oneObjFn()).CreateInBatches(data, 20).Error
	if err != nil {
		err = errors.Wrap(err, "db.CreateInBatches")
		return
	}
	return
}
