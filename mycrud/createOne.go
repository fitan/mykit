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

func (c *CRUD) createOneHandler() {
	c.handler(CreateOneMethodName, http.MethodPost, "/{tableName}", c.createOneEndpoint(), c.createOneDecode())
}

type CreateOneRequest struct {
	TableName string `json:"tableName"`
	Body      interface{}
}

func (c *CRUD) createOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]

		msg, err := c.tableMsg(req.TableName)
		if err != nil {
			return
		}

		body := msg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return nil, err
		}

		req.Body = body
		return req, err
	}
}

func (c *CRUD) createOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateOneRequest)
		err = c.createOne(ctx, req.TableName, req.Body)
		return nil, err
	}
}

func (c *CRUD) createOne(ctx context.Context, tableName string, data interface{}) (err error) {
	_, err = c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Table(tableName).Create(&data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Table(tableName).Create(&data).Error")
		return
	}
	return
}
