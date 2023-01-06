package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

type UpdateOneImpl interface {
	UpdateOneHandler()
	UpdateOneDecode() kithttp.DecodeRequestFunc
	UpdateOneEndpoint() endpoint.Endpoint
	UpdateOne(ctx context.Context, id string, data interface{}) (err error)
}

type UpdateOne struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (u *UpdateOne) UpdateOneHandler() {
	u.Crud.Handler(UpdateOneMethodName, http.MethodPut, "/"+u.TableMsg.schema.Table+"/{id}", u.UpdateOneEndpoint(), u.UpdateOneDecode())
}

type UpdateOneRequest struct {
	Id   string      `json:"id"`
	Body interface{} `json:"body"`
}

func (u *UpdateOne) UpdateOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := UpdateOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]

		body := u.TableMsg.oneObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			err = errors.Wrap(err, "json.NewDecoder(r.Body).Decode(body)")
			return
		}

		req.Body = body
		return req, err

	}
}

func (u *UpdateOne) UpdateOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateOneRequest)
		err = u.UpdateOne(ctx, req.Id, req.Body)
		return nil, err
	}
}

func (u *UpdateOne) UpdateOne(ctx context.Context, id string, data interface{}) (err error) {

	db, commit := u.Crud.db.Tx(ctx)
	defer commit(err)

	err = db.Model(u.TableMsg.oneObjFn()).Select("*").Where("id = ?", id).Updates(data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Updates()")
		return
	}
	return
}
