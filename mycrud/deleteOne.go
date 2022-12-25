package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func (c *CRUD) deleteOneHandler() {
	c.handler(DeleteOneMethodName, http.MethodDelete, "/{tableName}/{id}", c.deleteOneEndpoint(), c.deleteOneDecode())
}

type deleteOneRequest struct {
	TableName string `json:"tableName"`
	Id        string `json:"id"`
}

func (c *CRUD) deleteOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := deleteOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		return req, nil
	}
}

func (c *CRUD) deleteOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteOneRequest)
		err = c.deleteOne(ctx, req.TableName, req.Id)
		return nil, err
	}
}

func (c *CRUD) deleteOne(ctx context.Context, tableName, id string) (err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Table(tableName).Where("id = ?", id).Delete(&(msg.oneObjFn)).Error
	return
}
