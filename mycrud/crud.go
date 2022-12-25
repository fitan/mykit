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
	GetOneMethodName     = "GetOne"
	GetManyMethodName    = "GetMany"
	CreateOneMethodName  = "CreateOne"
	CreateManyMethodName = "CreateMany"
	UpdateOneMethodName  = "UpdateOne"
	UpdateManyMethodName = "UpdateMany"
	DeleteOneMethodName  = "DeleteOne"
	DeleteManyMethodName = "DeleteMany"
)

type CRUD struct {
	tables  map[string]tableMsg
	m       *mux.Router
	db      *mygorm.DB
	enc     kithttp.EncodeResponseFunc
	options []kithttp.ServerOption
}

type methodMsg struct {
	getOneHas     bool
	getOne        GetOneActionMethod
	getManyHas    bool
	getMany       GetManyActionMethod
	updateOneHas  bool
	updateOne     UpdateOneActionMethod
	updateManyHas bool
	updateMany    UpdateManyActionMethod
	deleteOneHas  bool
	deleteOne     DeleteOneActionMethod
	deleteManyHas bool
	deleteMany    DeleteManyActionMethod
	createOneHas  bool
	createOne     CreateOneActionMethod
	createManyHas bool
	createMany    CreateManyActionMethod
	enc           kithttp.EncodeResponseFunc
	options       []kithttp.ServerOption
}

func NewCRUD(m *mux.Router, db *gorm.DB, encode kithttp.EncodeResponseFunc) *CRUD {
	enc := kithttp.EncodeJSONResponse
	if encode != nil {
		enc = encode
	}
	return &CRUD{m: m, tables: map[string]tableMsg{}, db: mygorm.New(db), enc: enc}
}

func (c *CRUD) RegisterTable(oneObjFn func() interface{}, manyObjFn func() interface{}) error {
	tSchema, err := schema.Parse(oneObjFn(), &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return errors.Wrap(err, "schema.Parse")
	}

	c.tables[tSchema.Table] = tableMsg{
		oneObjFn:  oneObjFn,
		manyObjFn: manyObjFn,
		schema:    *tSchema,
	}
	return nil
}

type tableMsg struct {
	oneObjFn   func() interface{}
	manyObjFn  func() interface{}
	schema     schema.Schema
	encMap     map[string]kithttp.EncodeResponseFunc
	optionsMap map[string][]kithttp.ServerOption
	methodMap  map[string]methodMsg
}

func (c *CRUD) RegisterMethod(tableName string) *RegisterMethodHelper {
	return &RegisterMethodHelper{crud: c, tableName: tableName}
}

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
	c.getOneHandler()
	c.getManyHandler()
	c.updateOneHandler()
	c.updateManyHandler()
	c.createOneHandler()
	c.createManyHandler()
	c.deleteOneHandler()
	c.deleteManyHandler()
}

func (c *CRUD) tableMsg(tableName string) (tableMsg, error) {
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

		return mygorm.SetScopesToCtx(ctx, request, msg.schema)
	})
}

func (c *CRUD) handler(methodName, httpMethod string, path string, e endpoint.Endpoint, dec kithttp.DecodeRequestFunc, opts ...kithttp.ServerOption) {
	c.m.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		tableName := mux.Vars(request)["tableName"]
		msg, err := c.tableMsg(tableName)
		if err != nil {
			myhttp.ResponseJsonEncode(writer, map[string]interface{}{"err": err.Error()})
			return
		}
		enc := msg.encMap[methodName]
		if enc == nil {
			enc = c.enc
		}

		o := append(c.options, msg.optionsMap[methodName]...)
		o = append(o, opts...)

		kithttp.NewServer(e, dec, enc, o...).ServeHTTP(writer, request)
	}).Methods(httpMethod)

}
