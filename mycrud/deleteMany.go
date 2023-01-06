package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"strings"
)

type DeleteManyImpl interface {
	DeleteManyHandler()
	DeleteManyDecode() kithttp.DecodeRequestFunc
	DeleteManyEndpoint() endpoint.Endpoint
	DeleteMany(ctx context.Context, ids []string) (err error)
}

type DeleteManyRequest struct {
	Ids []string `json:"ids"`
}

type DeleteMany struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (d *DeleteMany) DeleteManyHandler() {
	d.Crud.Handler(DeleteManyMethodName, http.MethodDelete, "/{tableName}", d.DeleteManyEndpoint(), d.DeleteManyDecode())
}

func (d *DeleteMany) DeleteManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := DeleteManyRequest{}
		ids := r.URL.Query().Get("ids")
		req.Ids = strings.Split(ids, ",")
		return req, nil
	}
}

func (d *DeleteMany) DeleteManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteManyRequest)
		err = d.DeleteMany(ctx, req.Ids)
		return nil, err
	}
}

func (d *DeleteMany) DeleteMany(ctx context.Context, ids []string) (err error) {

	db, commit := d.Crud.db.Tx(ctx)
	defer commit(err)

	err = db.Model(d.TableMsg.oneObjFn()).Where("id in (?)", ids).Delete(d.TableMsg.oneObjFn()).Error
	return
}
