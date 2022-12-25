package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func (c *CRUD) getOneHandler() {
	c.handler(GetOneMethodName, http.MethodGet, "/{tableName}/{id}", c.getOneEndpoint(), c.getOneDecode())
}

type GetOneRequest struct {
	TableName string `json:"table_name"`
	Id        string `json:"id"`
}

func (c *CRUD) getOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		return req, nil
	}
}

func (c *CRUD) getOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOneRequest)
		res, err := c.getOne(ctx, req.TableName, req.Id)
		return res, err
	}
}

func (c *CRUD) getOne(ctx context.Context, tableName, id string) (data interface{}, err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	db := c.db.Db(ctx)

	data = msg.oneObjFn()
	err = db.Table(tableName).First(&data).Error
	return
}
