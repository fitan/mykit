package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type DeleteManyRequest struct {
	TableName string   `json:"tableName"`
	Ids       []string `json:"ids"`
}

func (c *CRUD) DeleteManyHandler() {
	c.Handler(DeleteManyMethodName, http.MethodDelete, "/{tableName}", c.DeleteManyEndpoint(), c.DeleteManyDecode())
}

func (c *CRUD) DeleteManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := DeleteManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		ids := r.URL.Query().Get("ids")
		req.Ids = strings.Split(ids, ",")
		return req, nil
	}
}

func (c *CRUD) DeleteManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteManyRequest)
		err = c.DeleteMany(ctx, req.TableName, req.Ids)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) DeleteMany(ctx context.Context, tableName string, ids []string) (err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
	defer commit(err)

	err = db.Table(tableName).Where("id in (?)", ids).Delete(&(msg.oneObjFn)).Error
	return
}
