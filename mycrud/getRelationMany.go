package mycrud

import (
	"context"
	"fmt"
	"github.com/fitan/mykit/mygorm"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
)

func (c *CRUD) getRelationManyHandler() {
	c.handler(GetRelationManyMethodName, http.MethodGet, "/{tableName}/{id}/{relationTableName}/many", c.getRelationManyEndpoint(), c.getRelationManyDecode(), c.KitGormScopesBefore())
}

type GetRelationManyRequest struct {
	TableName         string `json:"tableName"`
	Id                string `json:"id"`
	RelationTableName string `json:"RelationTableName"`
}

func (c *CRUD) getRelationManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetRelationManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		return req, nil
	}
}

func (c *CRUD) getRelationManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRelationManyRequest)
		res, err := c.getRelationMany(ctx, req.TableName, req.Id, req.RelationTableName, nil)
		return c.endpointWrap(res, err)
	}
}

func (c *CRUD) getRelationMany(ctx context.Context, tableName, id, relationTableName string, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
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

	var total int64
	totalDB := c.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Model(relationTableNameMsg.oneObjFn()).Where(relationPrimaryKey+" in (?)", totalDB.Session(&gorm.Session{NewDB: true}).Table(tableName).Select(relationForeignKey).Where(relationForeignKey+" = ?", id)).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Count")
		return
	}

	db := c.db.Db(ctx).Model(relationTableNameMsg.oneObjFn()).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}

	list := relationTableNameMsg.manyObjFn()

	db.Where(relationPrimaryKey+" in (?)", db.Session(&gorm.Session{NewDB: true}).Table(tableName).Select(relationForeignKey).Where(relationForeignKey+" = ?", id)).Find(list)

	data.List = list
	data.Total = total

	return data, nil
}