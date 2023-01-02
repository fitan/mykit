package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func (c *CRUD) DeleteOneHandler() {
	c.Handler(DeleteOneMethodName, http.MethodDelete, "/{tableName}/{id}", c.DeleteOneEndpoint(), c.DeleteOneDecode())
}

type DeleteOneRequest struct {
	TableName string `json:"tableName"`
	Id        string `json:"id"`
}

func (c *CRUD) DeleteOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := DeleteOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		return req, nil
	}
}

func (c *CRUD) DeleteOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteOneRequest)
		err = c.DeleteOne(ctx, req.TableName, req.Id)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) DeleteOne(ctx context.Context, tableName, id string) (err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Model(msg.oneObjFn()).Where("id = ?", id).Delete(msg.oneObjFn()).Error
	return
}
