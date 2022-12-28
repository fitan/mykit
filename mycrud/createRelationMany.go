package mycrud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strconv"
)

func (c *CRUD) createRelationManyHandler() {
	c.handler(CreateRelationManyMethodName, http.MethodPost, "/{tableName}/{id}/{relationTableName}", c.createRelationManyEndpoint(), c.createRelationManyDecode())
}

type CreatRelationManyRequest struct {
	TableName         string      `json:"tableName"`
	Id                string      `json:"id"`
	RelationTableName string      `json:"relationName"`
	Body              interface{} `json:"body"`
}

func (c *CRUD) createRelationManyDecode() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		req := CreatRelationManyRequest{}
		v := mux.Vars(r)
		req.TableName = v["tableName"]
		req.Id = v["id"]
		req.RelationTableName = v["relationTableName"]
		msg, err := c.tableMsg(req.RelationTableName)
		if err != nil {
			return
		}
		req.Body = msg.manyObjFn()
		err = json.NewDecoder(r.Body).Decode(&req.Body)
		if err != nil {
			err = errors.Wrap(err, "json.Decode")
			return
		}
		return req, nil
	}
}

func (c *CRUD) createRelationManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreatRelationManyRequest)
		res, err := c.createRelationMany(ctx, req.TableName, req.Id, req.RelationTableName, req.Body)
		return c.endpointWrap(res, err)
	}
}

func (c *CRUD) createRelationMany(ctx context.Context, tableName, id, relationTableName string, body interface{}) (data interface{}, err error) {
	msg, err := c.tableMsg(tableName)
	if err != nil {
		return
	}

	var hasRelation bool
	var relationFieldName string

	for k, v := range msg.schema.Relationships.Relations {
		if v.FieldSchema.Table == relationTableName {
			hasRelation = true
			relationFieldName = k
		}
	}

	if !hasRelation {
		err = fmt.Errorf("table %s has no relation with table %s", tableName, relationTableName)
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

	err = c.db.Db(ctx).Model(model).Association(relationFieldName).Append(body)
	if err != nil {
		err = errors.Wrap(err, "db.Append")
		return
	}

	return nil, nil
}
