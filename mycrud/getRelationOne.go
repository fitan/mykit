package mycrud

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
)

func (c *CRUD) getRelationOneHandler() {
	c.handler(GetRelationOneMethodName, http.MethodGet, "/{tableName}/{id}/{relationTableName}", c.getRelationOneEndpoint(), c.getRelationOneDecode())
}

type GetRelationOneRequest struct {
	TableName         string `json:"tableName"`
	Id                string `json:"id"`
	RelationTableName string `json:"RelationTableName"`
}

func (c *CRUD) getRelationOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetRelationOneRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		return req, nil
	}
}

func (c *CRUD) getRelationOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRelationOneRequest)
		res, err := c.getRelationOne(ctx, req.TableName, req.Id, req.RelationTableName, nil)
		return c.endpointWrap(res, err)
	}
}

func (c *CRUD) getRelationOne(ctx context.Context, tableName, id, relationTableName string, scopes []func(db *gorm.DB) *gorm.DB) (data interface{}, err error) {
	relationTableNameMsg, err := c.tableMsg(relationTableName)
	if err != nil {
		return
	}

	var _ tableMsg
	var relation schema.Relationship
	var hasRelation bool

	for _, v := range relationTableNameMsg.schema.Relationships.Relations {
		if v.FieldSchema.Table == tableName {
			relation = *v
			_, err = c.tableMsg(relationTableName)
			if err != nil {
				return
			}
			hasRelation = true
			break
		}
	}

	if hasRelation == false {
		err = errors.New("no relation")
		return
	}

	if len(relation.References) == 0 {
		err = fmt.Errorf("not found reference: %s", tableName)
		return
	}

	relationForeignKey := relation.References[0].ForeignKey.DBName
	relationPrimaryKey := relation.References[0].PrimaryKey.DBName

	data = relationTableNameMsg.oneObjFn()

	db := c.db.Db(ctx)
	db.Model(relationTableNameMsg.oneObjFn()).Where(relationPrimaryKey+" in (?)", db.Session(&gorm.Session{NewDB: true}).Table(tableName).Select(relationForeignKey).Where(relationForeignKey+" = ?", id)).First(data)

	return data, nil
}
