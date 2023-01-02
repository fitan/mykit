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

func (c *CRUD) CreateOneHandler() {
	c.Handler(CreateOneMethodName, http.MethodPost, "/{tableName}", c.CreateOneEndpoint(), c.CreateOneDecode())
}

type CreateOneRequest struct {
	TableName string `json:"tableName"`
	Body      interface{}
}

func (c *CRUD) CreateOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]

		msg, err := c.tableMsg(req.TableName)
		if err != nil {
			return
		}

		body := msg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return nil, err
		}

		req.Body = body
		return req, err
	}
}

func (c *CRUD) CreateOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateOneRequest)
		err = c.CreateOne(ctx, req.TableName, req.Body)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) CreateOne(ctx context.Context, tableName string, data interface{}) (err error) {
	_, err = c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Table(tableName).Create(data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Table(tableName).Create(data).Error")
		return
	}
	return
}
