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

func (c *CRUD) CreateManyHandler() {
	c.Handler(CreateManyMethodName, http.MethodPost, "/{tableName}/many", c.CreateManyEndpoint(), c.CreateManyDecode())
}

type CreateManyRequest struct {
	TableName string `json:"table_name"`
	Body      interface{}
}

func (c *CRUD) CreateManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		msg, err := c.tableMsg(req.TableName)
		if err != nil {
			return
		}

		body := msg.manyObjFn()
		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return
		}

		req.Body = body
		return req, nil
	}
}

func (c *CRUD) CreateManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateManyRequest)
		err = c.CreateMany(ctx, req.TableName, req.Body)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) CreateMany(ctx context.Context, tableName string, data interface{}) (err error) {
	_, err = c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Table(tableName).CreateInBatches(data, 20).Error
	if err != nil {
		err = errors.Wrap(err, "db.Table(tableName).Create(data).Error")
		return
	}
	return
}
