package mycrud

import (
	"context"
	"fmt"
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
	GetOneMethodName             = "GetOne"
	GetManyMethodName            = "GetMany"
	CreateOneMethodName          = "CreateOne"
	CreateManyMethodName         = "CreateMany"
	UpdateOneMethodName          = "UpdateOne"
	UpdateManyMethodName         = "UpdateMany"
	DeleteOneMethodName          = "DeleteOne"
	DeleteManyMethodName         = "DeleteMany"
	RelationGetOneMethodName     = "RelationGetOne"
	RelationGetManyMethodName    = "RelationGetMany"
	RelationCreateOneMethodName  = "RelationCreateOne"
	RelationCreateManyMethodName = "RelationCreateMany"
)

type Core struct {
	tables      map[string]*TableMsg
	handlers    map[string][]func(core *Core, msg *TableMsg)
	m           *mux.Router
	db          *mygorm.DB
	enc         kithttp.EncodeResponseFunc
	options     []kithttp.ServerOption
	endpointMid []endpoint.Middleware
	permissions Permissions
}

type Permissions func(ctx context.Context, tableName string, methodName string) (bool, error)

func NewCore(m *mux.Router, db *gorm.DB, encode kithttp.EncodeResponseFunc, opts []kithttp.ServerOption) *Core {
	enc := myhttp.EncodeJSONResponse
	if encode != nil {
		enc = encode
	}
	core := &Core{m: m, tables: map[string]*TableMsg{}, handlers: map[string][]func(core *Core, msg *TableMsg){}, db: mygorm.New(db), enc: enc, options: make([]kithttp.ServerOption, 0), endpointMid: make([]endpoint.Middleware, 0)}
	core.options = append(core.options, myhttp.KitErrorEncoder())
	core.options = append(core.options, opts...)
	return core
}

func (c *Core) kitDtoEncodeJsonResponse(dto func(i interface{}) interface{}) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		if dto != nil {
			i = dto(i)
		}
		return c.enc(ctx, writer, i)
	}
}

func (c *Core) RegisterCRUD(oneObjFn func() interface{}, manyObjFn func() interface{}, crud CRUD) (*TableMsg, error) {
	tSchema, err := schema.Parse(oneObjFn(), &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, errors.Wrap(err, "schema.Parse")
	}

	t := &TableMsg{
		oneObjFn:  oneObjFn,
		manyObjFn: manyObjFn,
		schema:    *tSchema,
	}

	c.tables[tSchema.Table] = t
	if crud.GetOne != nil {
		c.RegHandler(crud.GetOne)
	}
	if crud.GetMany != nil {
		c.RegHandler(crud.GetMany)
	}
	if crud.CreateOne != nil {
		c.RegHandler(crud.CreateOne)
	}
	if crud.CreateMany != nil {
		c.RegHandler(crud.CreateMany)
	}
	if crud.UpdateOne != nil {
		c.RegHandler(crud.UpdateOne)
	}
	if crud.UpdateMany != nil {
		c.RegHandler(crud.UpdateMany)
	}
	if crud.DeleteOne != nil {
		c.RegHandler(crud.DeleteOne)
	}
	if crud.DeleteMany != nil {
		c.RegHandler(crud.DeleteMany)
	}
	for _, v := range crud.Relations {
		if v.GetOne != nil {
			c.RegHandler(v.GetOne)
		}
		if v.GetMany != nil {
			c.RegHandler(v.GetMany)
		}
		if v.CreateOne != nil {
			c.RegHandler(v.CreateOne)
		}
		if v.CreateMany != nil {
			c.RegHandler(v.CreateMany)
		}
	}

	for _, v := range crud.Methods {
		c.RegHandler(v)
	}
	return c.tables[tSchema.Table], nil

}

func (c *Core) RegisterTable(oneObjFn func() interface{}, manyObjFn func() interface{}, regs ...func(core *Core, tableMsg *TableMsg)) (*TableMsg, error) {
	tSchema, err := schema.Parse(oneObjFn(), &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, errors.Wrap(err, "schema.Parse")
	}
	t := &TableMsg{
		oneObjFn:  oneObjFn,
		manyObjFn: manyObjFn,
		schema:    *tSchema,
	}

	c.tables[tSchema.Table] = t
	c.handlers[tSchema.Table] = regs
	return c.tables[tSchema.Table], nil
}

func (c *Core) Run() {
	for tableName, tableMsg := range c.tables {
		for _, reg := range c.handlers[tableName] {
			reg(c, tableMsg)
		}
	}
}

func (c *Core) GetTableMsg(tableName string) (*TableMsg, error) {
	msg, ok := c.tables[tableName]
	if !ok {
		return msg, fmt.Errorf("not found table %s", tableName)
	}
	return msg, nil
}

func (c *Core) KitGormScopesBefore(tableMsg *TableMsg) kithttp.ServerOption {
	return kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
		return mygorm.SetScopesToCtx(ctx, request, tableMsg.schema)
	})
}

func (c *Core) RegHandler(impl KitHttpImpl) {
	e := impl.GetEndpoint()

	for _, mid := range impl.GetEndpointMid() {
		e = mid(e)
	}

	for _, mid := range c.endpointMid {
		e = mid(e)
	}

	enc := impl.GetEncode()
	if enc == nil {
		enc = c.enc
	}

	opts := c.options
	opts = append(opts, impl.GetOptions()...)
	c.m.Name(impl.GetName()).Methods(impl.GetHttpMethod()).Path(impl.GetHttpPath()).Handler(kithttp.NewServer(
		e,
		impl.GetDecode(),
		enc,
		opts...,
	))
}

func NewCrud(core *Core, tableMsg *TableMsg) {
	core.RegHandler(NewGetOne(core, tableMsg))
	core.RegHandler(NewGetMany(core, tableMsg))

	core.RegHandler(NewUpdateOne(core, tableMsg))
	core.RegHandler(NewUpdateMany(core, tableMsg))

	core.RegHandler(NewCreateOne(core, tableMsg))
	core.RegHandler(NewCreateMany(core, tableMsg))

	core.RegHandler(NewDeleteOne(core, tableMsg))
	core.RegHandler(NewDeleteMany(core, tableMsg))

	for _, impl := range newGetRelationOne(core, tableMsg) {
		core.RegHandler(impl)
	}

	for _, impl := range newGetRelationMany(core, tableMsg) {
		core.RegHandler(impl)
	}

	for _, impl := range newCreateRelationOne(core, tableMsg) {
		core.RegHandler(impl)
	}

	for _, impl := range newCreateRelationMany(core, tableMsg) {
		core.RegHandler(impl)
	}
}

func NewRepo(core *Core, msg *TableMsg) *Repo {
	return &Repo{
		Core:     core,
		TableMsg: msg,
	}
}

func NewGetOne(core *Core, tableMsg *TableMsg) *GetOne {
	return &GetOne{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + GetOneMethodName,
			HttpMethod: http.MethodGet,
			HttpPath:   fmt.Sprintf("/%s/{id}", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func NewGetMany(core *Core, tableMsg *TableMsg) *GetMany {
	return &GetMany{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + GetManyMethodName,
			HttpMethod: http.MethodGet,
			HttpPath:   fmt.Sprintf("/%s", tableMsg.schema.Table),
			Encode:     nil,
			Options:    []kithttp.ServerOption{core.KitGormScopesBefore(tableMsg)},
		},
	}
}

func NewCreateOne(core *Core, tableMsg *TableMsg) *CreateOne {
	return &CreateOne{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + CreateOneMethodName,
			HttpMethod: http.MethodPost,
			HttpPath:   fmt.Sprintf("/%s", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func NewCreateMany(core *Core, tableMsg *TableMsg) *CreateMany {
	return &CreateMany{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + CreateManyMethodName,
			HttpMethod: http.MethodPost,
			HttpPath:   fmt.Sprintf("/%s/batch", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func NewUpdateOne(core *Core, tableMsg *TableMsg) *UpdateOne {
	return &UpdateOne{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + UpdateOneMethodName,
			HttpMethod: http.MethodPut,
			HttpPath:   fmt.Sprintf("/%s/{id}", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func NewUpdateMany(core *Core, tableMsg *TableMsg) *UpdateMany {
	return &UpdateMany{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + UpdateManyMethodName,
			HttpMethod: http.MethodPut,
			HttpPath:   fmt.Sprintf("/%s/batch", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func NewDeleteOne(core *Core, tableMsg *TableMsg) *DeleteOne {
	return &DeleteOne{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + DeleteOneMethodName,
			HttpMethod: http.MethodDelete,
			HttpPath:   fmt.Sprintf("/%s/{id}", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func NewDeleteMany(core *Core, tableMsg *TableMsg) *DeleteMany {
	return &DeleteMany{
		Repo: NewRepo(core, tableMsg),
		KitHttpConfig: &KitHttpConfig{
			Name:       tableMsg.schema.Table + DeleteManyMethodName,
			HttpMethod: http.MethodDelete,
			HttpPath:   fmt.Sprintf("/%s/batch", tableMsg.schema.Table),
			Encode:     nil,
			Options:    nil,
		},
	}
}

func newGetRelationOne(core *Core, msg *TableMsg) (res []*GetRelationOne) {
	for _, v := range msg.schema.Relationships.HasOne {
		relationTableName := v.FieldSchema.Table

		res = append(res, &GetRelationOne{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + GetOneMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodGet,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
				Encode:     nil,
				Options:    nil,
			},
			RelationTableName: relationTableName,
		})
	}

	for _, v := range msg.schema.Relationships.BelongsTo {
		relationTableName := v.FieldSchema.Table

		res = append(res, &GetRelationOne{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + GetOneMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodGet,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
				Encode:     nil,
				Options:    nil,
			},
			RelationTableName: relationTableName,
		})
	}
	return
}

func newGetRelationMany(core *Core, msg *TableMsg) (res []*GetRelationMany) {
	for _, v := range msg.schema.Relationships.HasMany {
		relationTableName := v.FieldSchema.Table
		relationTable, err := core.GetTableMsg(relationTableName)
		if err != nil {
			panic(err)
		}

		res = append(res, &GetRelationMany{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + GetManyMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodGet,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
				Encode:     nil,
				Options:    []kithttp.ServerOption{core.KitGormScopesBefore(relationTable)},
			},
			RelationTableName: relationTableName,
		})
	}

	for _, v := range msg.schema.Relationships.Many2Many {
		relationTableName := v.FieldSchema.Table
		relationTable, err := core.GetTableMsg(relationTableName)
		if err != nil {
			panic(err)
		}

		res = append(res, &GetRelationMany{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + GetManyMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodGet,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
				Encode:     nil,
				Options:    []kithttp.ServerOption{core.KitGormScopesBefore(relationTable)},
			},
			RelationTableName: relationTableName,
		})
	}
	return
}

func newCreateRelationOne(core *Core, msg *TableMsg) (res []*CreateRelationOne) {
	for _, v := range msg.schema.Relationships.HasOne {
		relationTableName := v.FieldSchema.Table

		res = append(res, &CreateRelationOne{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + CreateOneMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodPost,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
				Encode:     nil,
				Options:    nil,
			},
			RelationTableName: relationTableName,
		})
	}

	for _, v := range msg.schema.Relationships.BelongsTo {
		relationTableName := v.FieldSchema.Table

		res = append(res, &CreateRelationOne{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + CreateOneMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodPost,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
				Encode:     nil,
				Options:    nil,
			},
			RelationTableName: relationTableName,
		})
	}
	return
}

func newCreateRelationMany(core *Core, msg *TableMsg) (res []*CreateRelationMany) {
	for _, v := range msg.schema.Relationships.HasMany {
		relationTableName := v.FieldSchema.Table

		res = append(res, &CreateRelationMany{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + CreateManyMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodPost,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
			},
			RelationTableName: relationTableName,
		})
	}

	for _, v := range msg.schema.Relationships.Many2Many {
		relationTableName := v.FieldSchema.Table

		res = append(res, &CreateRelationMany{
			Repo: NewRepo(core, msg),
			KitHttpConfig: &KitHttpConfig{
				Name:       relationTableName + CreateManyMethodName + "By" + msg.schema.Table,
				HttpMethod: http.MethodPost,
				HttpPath:   fmt.Sprintf("/%s/{id}/%s", msg.schema.Table, relationTableName),
			},
			RelationTableName: relationTableName,
		})
	}
	return
}
