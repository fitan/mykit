package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type DeleteOneImpl interface {
	DeleteOneHandler()
	DeleteOneDecode() kithttp.DecodeRequestFunc
	DeleteOneEndpoint() endpoint.Endpoint
	DeleteOne(ctx context.Context, id string) (err error)
}

type DeleteOne struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (d *DeleteOne) DeleteOneHandler() {
	d.Crud.Handler(DeleteOneMethodName, http.MethodDelete, "/"+d.TableMsg.schema.Table+"/{id}", d.DeleteOneEndpoint(), d.DeleteOneDecode())
}

type DeleteOneRequest struct {
	Id string `json:"id"`
}

func (d *DeleteOne) DeleteOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := DeleteOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		return req, nil
	}
}

func (d *DeleteOne) DeleteOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteOneRequest)
		err = d.DeleteOne(ctx, req.Id)
		return nil, err
	}
}

func (d *DeleteOne) DeleteOne(ctx context.Context, id string) (err error) {

	db, commit := d.Crud.db.Tx(ctx)
	defer commit(err)

	err = db.Model(d.TableMsg.oneObjFn()).Where("id = ?", id).Delete(d.TableMsg.oneObjFn()).Error
	return
}
