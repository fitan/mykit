package mycrud

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type GetOneImpl interface {
	GetOneHandler()
	GetOneDecode() kithttp.DecodeRequestFunc
	GetOneEndpoint() endpoint.Endpoint
	GetOne(ctx context.Context, tableName, id string) (data interface{}, err error)
}

type GetOne struct {
	Crud       *Core
	TableMsg   *tableMsg
	MethodName string
	HttpMethod string
	HttpPath   string
	Serializer func(i interface{}) interface{}
}

func (g *GetOne) GetOneHandler() {
	g.Crud.Handler(GetOneMethodName, http.MethodGet, "/"+g.TableMsg.schema.Table+"/{id}", g.GetOneEndpoint(), g.GetOneDecode())
}

type GetOneRequest struct {
	Id string `json:"id"`
}

func (g *GetOne) GetOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		return req, nil
	}
}

func (g *GetOne) GetOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOneRequest)
		res, err := g.GetOne(ctx, req.Id)
		return res, err
	}
}

func (g *GetOne) GetOne(ctx context.Context, id string) (data interface{}, err error) {

	db := g.Crud.db.Db(ctx)

	obj := g.TableMsg.oneObjFn()
	err = db.Model(g.TableMsg.oneObjFn()).Where("id = ?", id).First(obj).Error
	return obj, err
}
