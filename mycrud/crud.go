package mycrud

import (
	"context"
	"fmt"
	"github.com/fitan/mykit/myctx"
	"github.com/fitan/mykit/mygorm"
	"github.com/fitan/mykit/myhttp"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"sync"
)

const (
	GetOneMethodName         = "GetOne"
	GetManyMethodName        = "GetMany"
	CreateOneMethodName      = "CreateOne"
	CreateManyMethodName     = "CreateMany"
	UpdateOneMethodName      = "UpdateOne"
	UpdateManyMethodName     = "UpdateMany"
	DeleteOneMethodName      = "DeleteOne"
	DeleteManyMethodName     = "DeleteMany"
	MethodMethodName         = "Method"
	RelationMethodMethodName = "RelationMethod"
)

type CRUD struct {
	tables       map[string]*tableMsg
	m            *mux.Router
	db           *mygorm.DB
	enc          kithttp.EncodeResponseFunc
	endpointWrap func(response interface{}, err error) (interface{}, error)
	options      []kithttp.ServerOption
}

//type methodMsg struct {
//	getOneHas     bool
//	getOne        GetOneActionMethod
//	getManyHas    bool
//	GetMany       GetManyActionMethod
//	updateOneHas  bool
//	updateOne     UpdateOneActionMethod
//	updateManyHas bool
//	updateMany    UpdateManyActionMethod
//	deleteOneHas  bool
//	deleteOne     DeleteOneActionMethod
//	deleteManyHas bool
//	deleteMany    DeleteManyActionMethod
//	createOneHas  bool
//	createOne     CreateOneActionMethod
//	createManyHas bool
//	createMany    CreateManyActionMethod
//	enc           kithttp.EncodeResponseFunc
//	options       []kithttp.ServerOption
//}

func NewCRUD(m *mux.Router, db *gorm.DB, encode kithttp.EncodeResponseFunc, opts []kithttp.ServerOption) *CRUD {
	enc := myhttp.EncodeJSONResponse
	if encode != nil {
		enc = encode
	}
	crud := &CRUD{m: m, tables: map[string]*tableMsg{}, db: mygorm.New(db), enc: enc, endpointWrap: myhttp.WrapResponse, options: make([]kithttp.ServerOption, 0)}
	crud.options = append(crud.options, myhttp.KitErrorEncoder())
	crud.options = append(crud.options, opts...)
	return crud
}

func (c *CRUD) kitDtoEncodeJsonResponse(dto func(i interface{}) interface{}) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		if dto != nil {
			i = dto(i)
		}
		return c.enc(ctx, writer, i)
	}
}

func (c *CRUD) RegisterTable(oneObjFn func() interface{}, manyObjFn func() interface{}) (*tableMsg, error) {
	tSchema, err := schema.Parse(oneObjFn(), &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, errors.Wrap(err, "schema.Parse")
	}
	t := &tableMsg{
		oneObjFn:   oneObjFn,
		manyObjFn:  manyObjFn,
		schema:     *tSchema,
		encMap:     map[string]kithttp.EncodeResponseFunc{},
		dtoMap:     map[string]func(i interface{}) interface{}{},
		optionsMap: map[string][]kithttp.ServerOption{},
	}

	t.Option(GetManyMethodName, c.KitGormScopesBefore())

	c.tables[tSchema.Table] = t
	return c.tables[tSchema.Table], nil
}

//func (c *CRUD) RegisterMethod(tableName string) *RegisterMethodHelper {
//	return &RegisterMethodHelper{crud: c, tableName: tableName}
//}

//func (c *CRUD) runMethod() {
//	c.m.HandleFunc("/{tableName}/method/{methodName}", func(writer http.ResponseWriter, request *http.Request) {
//		v := mux.Vars(request)
//		tableName := v["tableName"]
//		methodName := v["methodName"]
//		msg, err := c.tableMsg(tableName)
//		if err != nil {
//			myhttp.ResponseJsonEncode(writer, map[string]interface{}{"err": err.Error()})
//			return
//		}
//
//		methodMsg,ok := msg.methodMap[methodName]
//		if !ok {
//			myhttp.ResponseJsonEncode(writer, map[string]interface{}{"err": fmt.Sprintf("not found method %s", methodName)})
//			return
//		}
//
//		if methodMsg.createOneHas {
//			kithttp.NewServer(c.createOneEndpoint(), c.createOneDecode(),
//				func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
//					methodMsg.createOne(ctx, )
//				})
//		}
//
//
//	}
//}

func (c *CRUD) run() {
	c.GetOneHandler()
	c.GetManyHandler()
	c.GetRelationOneHandler()
	c.GetRelationManyHandler()
	c.UpdateOneHandler()
	c.UpdateManyHandler()
	c.CreateOneHandler()
	c.CreateManyHandler()
	c.CreateRelationManyHandler()
	c.CreateRelationOneHandler()
	c.DeleteOneHandler()
	c.DeleteManyHandler()
}

func (c *CRUD) tableMsg(tableName string) (*tableMsg, error) {
	msg, ok := c.tables[tableName]
	if !ok {
		return msg, fmt.Errorf("not found table %s", tableName)
	}
	return msg, nil
}

func (c *CRUD) KitGormScopesBefore() kithttp.ServerOption {
	return kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {

		tableName, ok := mux.Vars(request)["tableName"]
		if !ok {
			return ctx
		}

		msg, err := c.tableMsg(tableName)
		if err != nil {
			return context.WithValue(ctx, myctx.CtxGormScopesKey, mygorm.CtxGormScopesValue{
				Err: err,
			})
		}

		relationTableName, ok := mux.Vars(request)["relationTableName"]
		if ok {
			var err error
			msg, err = c.tableMsg(relationTableName)
			if err != nil {
				return context.WithValue(ctx, myctx.CtxGormScopesKey, mygorm.CtxGormScopesValue{
					Err: err,
				})
			}
		}

		return mygorm.SetScopesToCtx(ctx, request, msg.schema)
	})
}

func (c *CRUD) Handler(methodName, httpMethod string, path string, e endpoint.Endpoint, dec kithttp.DecodeRequestFunc, opts ...kithttp.ServerOption) {
	c.m.HandleFunc("/crud"+path, func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		tableName := vars["tableName"]
		msg, err := c.tableMsg(tableName)
		if err != nil {
			myhttp.ResponseJsonEncode(writer, myhttp.Response{Error: err.Error(), Code: 500})
			return
		}

		relationTableName, ok := vars["relationTableName"]
		if ok {
			msg, err = c.tableMsg(relationTableName)
			if err != nil {
				myhttp.ResponseJsonEncode(writer, myhttp.Response{Error: err.Error(), Code: 500})
			}
		}

		enc := c.kitDtoEncodeJsonResponse(msg.dtoMap[methodName])

		o := append(c.options, msg.optionsMap[methodName]...)
		o = append(o, opts...)

		kithttp.NewServer(e, dec, enc, o...).ServeHTTP(writer, request)
	}).Methods(httpMethod).Name(methodName)

}
