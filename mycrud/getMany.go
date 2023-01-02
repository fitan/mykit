package mycrud

import (
	"context"
	"github.com/fitan/mykit/mygorm"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

func (c *CRUD) GetManyHandler() {
	c.Handler(GetManyMethodName, http.MethodGet, "/{tableName}", c.GetManyEndpoint(), c.GetManyDecode())
}

type GetManyData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type GetManyRequest struct {
	TableName string `json:"tableName"`
}

func (c *CRUD) GetManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]

		return req, nil
	}
}

func (c *CRUD) GetManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetManyRequest)
		res, err := c.GetMany(ctx, req.TableName, nil)
		return res, err
	}
}

func (c *CRUD) GetMany(ctx context.Context, tableName string, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	var total int64
	totalDB := c.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Model(msg.oneObjFn()).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Count")
		return
	}
	db := c.db.Db(ctx).Model(msg.oneObjFn()).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}
	//list := msg.oneObjFn()
	list := msg.manyObjFn()
	err = db.Find(list).Error
	if err != nil {
		err = errors.Wrap(err, "db.Find")
		return
	}

	data.Total = total
	data.List = list
	return
}
