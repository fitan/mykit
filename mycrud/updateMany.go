package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

type UpdateManyImpl interface {
	UpdateManyHandler()
	UpdateManyDecode() kithttp.DecodeRequestFunc
	UpdateManyEndpoint() endpoint.Endpoint
	UpdateMany(ctx context.Context, data interface{}) (err error)
}

type UpdateMany struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (u *UpdateMany) UpdateManyHandler() {
	u.Crud.Handler(UpdateManyMethodName, http.MethodPut, "/"+u.TableMsg.schema.Table, u.UpdateManyEndpoint(), u.UpdateManyDecode())
}

type UpdateManyRequest struct {
	Body interface{} `json:"body"`
}

func (u *UpdateMany) UpdateManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := UpdateManyRequest{}

		body := u.TableMsg.manyObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			err = errors.Wrap(err, "json.NewDecoder.Decode")
			return
		}

		req.Body = body
		return req, nil
	}
}

func (u *UpdateMany) UpdateManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateManyRequest)
		err = u.UpdateMany(ctx, req.Body)
		return nil, err
	}
}

func (u *UpdateMany) UpdateMany(ctx context.Context, data interface{}) (err error) {

	db, commit := u.Crud.db.Tx(ctx)
	defer commit(err)

	refV := reflect.ValueOf(data)

	if refV.Kind() == reflect.Ptr {
		refV = refV.Elem()
	}

	switch refV.Kind() {
	case reflect.Slice:
		for i := 0; i < refV.Len(); i++ {
			err = db.Model(refV.Index(i).Interface()).Updates(refV.Index(i).Interface()).Error
			if err != nil {
				err = errors.Wrap(err, "db.Updates()")
				return
			}
		}
	default:
		return errors.New("data must be slice")
	}

	return
}
