package mycrud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
)

type CreateRelationManyImpl interface {
	CreateRelationManyHandler()
	CreateRelationManyDecode() kithttp.DecodeRequestFunc
	CreateRelationManyEndpoint() endpoint.Endpoint
	CreateRelationMany(ctx context.Context, id, relationTableName string, body interface{}) (err error)
}

type CreateRelationMany struct {
	Crud     *Core
	TableMsg *tableMsg
}

func (c *CreateRelationMany) CreateRelationManyHandler() {
	c.Crud.Handler(CreateManyMethodName, http.MethodPost, "/"+c.TableMsg.schema.Table+"/{id}/{relationTableName}/many", c.CreateRelationManyEndpoint(), c.CreateRelationManyDecode())
}

type CreatRelationManyRequest struct {
	TableName         string      `json:"tableName"`
	Id                string      `json:"id"`
	RelationTableName string      `json:"relationTableName"`
	Body              interface{} `json:"body"`
}

func (c *CreateRelationMany) CreateRelationManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreatRelationManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		req.Body = c.TableMsg.manyObjFn()
		err = json.NewDecoder(r.Body).Decode(&req.Body)
		if err != nil {
			err = errors.Wrap(err, "json.Decode")
			return
		}
		return req, nil
	}
}

func (c *CreateRelationMany) CreateRelationManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreatRelationManyRequest)
		err = c.CreateRelationMany(ctx, req.Id, req.RelationTableName, req.Body)
		return nil, err
	}
}

func (c *CreateRelationMany) CreateRelationMany(ctx context.Context, id, relationTableName string, body interface{}) (err error) {
	db := c.Crud.db.Db(ctx)
	return CreateRelationManyService(ctx, db, c.tableMsg, id, relationTableName, body)
}

func CreateRelationManyService(ctx context.Context, db *gorm.DB, msg *tableMsg, id, relationTableName string, body interface{}) (err error) {

	var hasRelation bool
	var relationFieldName string

	for k, v := range msg.schema.Relationships.Relations {
		if v.FieldSchema.Table == relationTableName {
			hasRelation = true
			relationFieldName = k
		}
	}

	if !hasRelation {
		err = fmt.Errorf("table %s has no relation with table %s", msg.schema.Table, relationTableName)
		return
	}

	var gormID interface{}

	//Bool   DataType = "bool"
	//Int    DataType = "int"
	//Uint   DataType = "uint"
	//Float  DataType = "float"
	//String DataType = "string"
	//Time   DataType = "time"
	//Bytes  DataType = "bytes"
	switch msg.schema.FieldsByDBName["id"].FieldType.Kind().String() {
	case "int":
		gormID, err = strconv.Atoi(id)
		if err != nil {
			err = errors.Wrap(err, "strconv.Atoi")
			return
		}
	case "uint":
		gormID, err = strconv.Atoi(id)
		if err != nil {
			err = errors.Wrap(err, "strconv.Atoi")
			return
		}
		gormID = uint(gormID.(int))
	case "string":
		gormID = id
	default:
		err = fmt.Errorf("not support id type %s", msg.schema.FieldsByDBName["id"].DataType)
		return
	}

	model := msg.oneObjFn()
	reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(gormID))

	err = db.Model(model).Association(relationFieldName).Append(body)
	if err != nil {
		err = errors.Wrap(err, "db.Append")
		return
	}

	return
}
