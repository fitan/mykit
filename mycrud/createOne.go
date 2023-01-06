package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"net/http"
)

type CreateOneImpl interface {
	CreateOneHandler()
	CreateOneDecode() kithttp.DecodeRequestFunc
	CreateOneEndpoint() endpoint.Endpoint
	CreateOne(ctx context.Context, body interface{}) (err error)
}

type CreateOne struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (c *CreateOne) CreateOneHandler() {
	c.Crud.Handler(CreateOneMethodName, http.MethodPost, "/"+c.TableMsg.schema.Table, c.CreateOneEndpoint(), c.CreateOneDecode())
}

type CreateOneRequest struct {
	Body interface{}
}

func (c *CreateOne) CreateOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateOneRequest{}

		body := c.TableMsg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return nil, err
		}

		req.Body = body
		return req, err
	}
}

func (c *CreateOne) CreateOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateOneRequest)
		err = c.CreateOne(ctx, req.Body)
		return nil, err
	}
}

func (c *CreateOne) CreateOne(ctx context.Context, data interface{}) (err error) {

	db, commit := c.Crud.db.Tx(ctx)
	defer commit(err)

	err = db.Model(c.TableMsg.oneObjFn()).Create(data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Create")
		return
	}
	return
}
