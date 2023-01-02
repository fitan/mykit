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

func (c *CRUD) UpdateOneHandler() {
	c.Handler(UpdateOneMethodName, http.MethodPut, "/{tableName}/{id}", c.UpdateOneEndpoint(), c.UpdateOneDecode())
}

type UpdateOneRequest struct {
	TableName string      `json:"tableName"`
	Id        string      `json:"id"`
	Body      interface{} `json:"body"`
}

func (c *CRUD) UpdateOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := UpdateOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]

		msg, err := c.tableMsg(req.TableName)
		if err != nil {
			return
		}

		body := msg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			err = errors.Wrap(err, "json.NewDecoder(r.Body).Decode(body)")
			return
		}

		req.Body = body
		return req, err

	}
}

func (c *CRUD) UpdateOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateOneRequest)
		err = c.UpdateOne(ctx, req.TableName, req.Id, req.Body)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) UpdateOne(ctx context.Context, tableName string, id string, data interface{}) (err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Model(msg.oneObjFn()).Select("*").Where("id = ?", id).Updates(data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Updates()")
		return
	}
	return
}
