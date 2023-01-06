package mycrud

import (
	"context"
	"github.com/fitan/mykit/mygorm"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

type GetManyImpl interface {
	GetManyHandler()
	GetManyDecode() kithttp.DecodeRequestFunc
	GetManyEndpoint() endpoint.Endpoint
	GetMany(ctx context.Context, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error)
}

type GetMany struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (g *GetMany) GetManyHandler() {
	g.Crud.Handler(GetManyMethodName, http.MethodGet, "/"+g.TableMsg.schema.Table, g.GetManyEndpoint(), g.GetManyDecode())
}

type GetManyData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type GetManyRequest struct {
}

func (g *GetMany) GetManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		return nil, nil
	}
}

func (g *GetMany) GetManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := g.GetMany(ctx, nil)
		return res, err
	}
}

func (g *GetMany) GetMany(ctx context.Context, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
	var total int64
	totalDB := g.Crud.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Model(g.TableMsg.oneObjFn()).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Count")
		return
	}
	db := g.Crud.db.Db(ctx).Model(g.TableMsg.oneObjFn()).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}
	//list := msg.oneObjFn()
	list := g.TableMsg.manyObjFn()
	err = db.Find(list).Error
	if err != nil {
		err = errors.Wrap(err, "db.Find")
		return
	}

	data.Total = total
	data.List = list
	return
}
