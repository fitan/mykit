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

type GetRelationManyImpl interface {
	GetRelationManyHandler()
	GetRelationManyDecode() kithttp.DecodeRequestFunc
	GetRelationManyEndpoint() endpoint.Endpoint
	GetRelationMany(ctx context.Context, id string, data interface{}) (err error)
}

type GetRelationMany struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (g *GetRelationMany) GetRelationManyHandler() {
	g.Crud.Handler(GetManyMethodName, http.MethodGet, "/"+g.TableMsg.schema.Table+"/{id}/{relationTableName}/many", g.GetRelationManyEndpoint(), g.GetRelationManyDecode())
}

type GetRelationManyRequest struct {
	Id                string `json:"id"`
	RelationTableName string `json:"RelationTableName"`
}

func (g *GetRelationMany) GetRelationManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := GetRelationManyRequest{}
		v := mux.Vars(r)
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		return req, nil
	}
}

func (g *GetRelationMany) GetRelationManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRelationManyRequest)
		res, err := g.GetRelationMany(ctx, req.Id, req.RelationTableName, nil)
		return res, err
	}
}

func (g *GetRelationMany) GetRelationMany(ctx context.Context, id, relationTableName string, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
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

	var total int64
	totalDB := g.Crud.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Model(relationTableNameMsg.oneObjFn()).Where(relationForeignKey+" in (?)", totalDB.Session(&gorm.Session{NewDB: true}).Table(g.TableMsg.schema.Table).Select(relationPrimaryKey).Where(relationPrimaryKey+" = ?", id)).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Count")
		return
	}

	db := g.Crud.db.Db(ctx).Model(relationTableNameMsg.oneObjFn()).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}

	list := relationTableNameMsg.manyObjFn()

	db.Where(relationForeignKey+" in (?)", db.Session(&gorm.Session{NewDB: true}).Table(g.TableMsg.schema.Table).Select(relationPrimaryKey).Where(relationPrimaryKey+" = ?", id)).Find(list)

	data.List = list
	data.Total = total

	return data, nil
}
