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

func (c *CRUD) getManyHandler() {
	c.handler(GetManyMethodName, http.MethodGet, "/{tableName}", c.getManyEndpoint(), c.getManyDecode(), c.KitGormScopesBefore())
}

type GetManyData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type GetManyRequest struct {
	TableName string `json:"tableName"`
}

func (c *CRUD) getManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]

		return req, nil
	}
}

func (c *CRUD) getManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetManyRequest)
		res, err := c.getMany(ctx, req.TableName, nil)
		return res, err
	}
}

func (c *CRUD) getMany(ctx context.Context, tableName string, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	list := msg.manyObjFn()
	var total int64
	totalDB := c.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Table(tableName).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Table(tableName).Count(&total).Error")
		return
	}
	db := c.db.Db(ctx).Table(tableName).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}
	err = db.Find(&list).Error
	if err != nil {
		err = errors.Wrap(err, "db.Find(&list).Error")
		return
	}

	data.Total = total
	data.List = list
	return
}
