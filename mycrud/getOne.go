package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func (c *CRUD) GetOneHandler() {
	c.Handler(GetOneMethodName, http.MethodGet, "/{tableName}/{id}", c.GetOneEndpoint(), c.GetOneDecode())
}

type GetOneRequest struct {
	TableName string `json:"table_name"`
	Id        string `json:"id"`
}

func (c *CRUD) GetOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		return req, nil
	}
}

func (c *CRUD) GetOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOneRequest)
		res, err := c.GetOne(ctx, req.TableName, req.Id)
		return res, err
	}
}

func (c *CRUD) GetOne(ctx context.Context, tableName, id string) (data interface{}, err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	db := c.db.Db(ctx)

	obj := msg.oneObjFn()
	err = db.Model(msg.oneObjFn()).Where("id = ?", id).First(obj).Error
	return obj, err
}
