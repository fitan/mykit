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

type CreateRelationOneImpl interface {
	CreateRelationOneHandler()
	CreateRelationOneDecode() kithttp.DecodeRequestFunc
	CreateRelationOneEndpoint() endpoint.Endpoint
	CreateRelationOne(ctx context.Context, id, relationTableName string, body interface{}) (err error)
}

type CreateRelation struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (c *CreateRelation) CreateRelationOneHandler() {
	c.Crud.Handler(CreateOneMethodName, http.MethodPost, "/"+c.TableMsg.schema.Table+"/{id}/{relationTableName}", c.CreateRelationOneEndpoint(), c.CreateRelationOneDecode())
}

type CreateRelationOneRequest struct {
	Id                string      `json:"id"`
	RelationTableName string      `json:"relationTableName"`
	Body              interface{} `json:"body"`
}

func (c *CreateRelation) CreateRelationOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreateRelationOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		msg, err := c.Crud.tableMsg(req.RelationTableName)
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

func (c *CreateRelation) CreateRelationOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRelationOneRequest)
		err = c.CreateRelationOne(ctx, req.Id, req.RelationTableName, req.Body)
		return nil, err
	}
}

func (c *CreateRelation) CreateRelationOne(ctx context.Context, id, relationTableName string, body interface{}) (err error) {
	db := c.Crud.db.Db(ctx)
	return CreateRelationManyService(ctx, db, c.tableMsg, id, relationTableName, body)
}
