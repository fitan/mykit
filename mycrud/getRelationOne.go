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

type GetRelationOneImpl interface {
	GetRelationOneHandler()
	GetRelationOneDecode() kithttp.DecodeRequestFunc
	GetRelationOneEndpoint() endpoint.Endpoint
	GetRelationOne(ctx context.Context, id, relationTableName string, scopes []func(db *gorm.DB) *gorm.DB) (data interface{}, err error)
}

type GetRelationOne struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (g *GetRelationOne) GetRelationOneHandler() {
	g.Crud.Handler(GetOneMethodName, http.MethodGet, "/"+g.TableMsg.schema.Table+"/{id}/{relationTableName}", g.GetRelationOneEndpoint(), g.GetRelationOneDecode())
}

type GetRelationOneRequest struct {
	Id                string `json:"id"`
	RelationTableName string `json:"RelationTableName"`
}

func (g *GetRelationOne) GetRelationOneDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetRelationOneRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		return req, nil
	}
}

func (g *GetRelationOne) GetRelationOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRelationOneRequest)
		res, err := g.GetRelationOne(ctx, req.Id, req.RelationTableName, nil)
		return res, err
	}
}

func (g *GetRelationOne) GetRelationOne(ctx context.Context, id, relationTableName string, scopes []func(db *gorm.DB) *gorm.DB) (data interface{}, err error) {
	relationTableNameMsg, err := g.Crud.tableMsg(relationTableName)
	if err != nil {
		return
	}

	var _ tableMsg
	var relation schema.Relationship
	var hasRelation bool

	for _, v := range relationTableNameMsg.schema.Relationships.Relations {
		if v.FieldSchema.Table == g.TableMsg.schema.Table {
			relation = *v
			_, err = g.Crud.tableMsg(relationTableName)
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
		err = fmt.Errorf("not found reference: %s", g.TableMsg.schema.Table)
		return
	}

	relationForeignKey := relation.References[0].ForeignKey.DBName
	relationPrimaryKey := relation.References[0].PrimaryKey.DBName

	data = relationTableNameMsg.oneObjFn()

	db := g.Crud.db.Db(ctx)
	db.Model(relationTableNameMsg.oneObjFn()).Where(relationPrimaryKey+" in (?)", db.Session(&gorm.Session{NewDB: true}).Table(g.TableMsg.schema.Table).Select(relationForeignKey).Where(relationForeignKey+" = ?", id)).First(data)

	return data, nil
}
