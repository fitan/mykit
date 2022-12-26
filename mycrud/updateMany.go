package mycrud

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

func (c *CRUD) updateManyHandler() {
	c.handler(UpdateManyMethodName, http.MethodPut, "/{tableName}", c.updateManyEndpoint(), c.updateManyDecode())
}

type updateManyRequest struct {
	TableName string      `json:"tableName"`
	Body      interface{} `json:"body"`
}

func (c *CRUD) updateManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := updateManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]

		msg, err := c.tableMsg(req.TableName)
		if err != nil {
			return
		}

		body := msg.manyObjFn()

		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			err = errors.Wrap(err, "json.NewDecoder(r.Body).Decode(&body)")
			return
		}

		req.Body = body
		return req, nil
	}
}

func (c *CRUD) updateManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateManyRequest)
		err = c.updateMany(ctx, req.TableName, req.Body)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) updateMany(ctx context.Context, tableName string, data interface{}) (err error) {
	_, err = c.tableMsg(tableName)
	if err != nil {
		return
	}

	db, commit := c.db.Tx(ctx)
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
