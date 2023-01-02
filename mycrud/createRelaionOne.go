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

func (c *CRUD) CreateRelationOneHandler() {
	c.Handler(CreateOneMethodName, http.MethodPost, "/{tableName}/{id}/{relationTableName}", c.CreateRelationOneEndpoint(), c.CreateRelationOneDecode())
}

type CreateRelationOneRequest struct {
	TableName         string      `json:"tableName"`
	Id                string      `json:"id"`
	RelationTableName string      `json:"relationTableName"`
	Body              interface{} `json:"body"`
}

func (c *CRUD) CreateRelationOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateRelationOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		msg, err := c.tableMsg(req.RelationTableName)
		if err != nil {
			return
		}
		req.Body = msg.oneObjFn()
		err = json.NewDecoder(r.Body).Decode(&req.Body)
		if err != nil {
			err = errors.Wrap(err, "json.Decode")
			return
		}
		return req, nil
	}
}

func (c *CRUD) CreateRelationOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRelationOneRequest)
		err = c.CreateRelationOne(ctx, req.TableName, req.Id, req.RelationTableName, req.Body)
		return c.endpointWrap(nil, err)
	}
}

func (c *CRUD) CreateRelationOne(ctx context.Context, tableName, id, relationTableName string, body interface{}) (err error) {
	return c.CreateRelationMany(ctx, tableName, id, relationTableName, body)
}
