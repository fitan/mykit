package mygorm

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/fitan/mykit/myctx"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CtxGormScopesValue struct {
	QScopes     []func(db *gorm.DB) *gorm.DB
	OtherScopes []func(db *gorm.DB) *gorm.DB
	Err         error
}

func KitScopesBefore(i interface{}) (option kithttp.ServerOption, err error) {

	tSchema, err := schema.Parse(i, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		err = errors.Wrap(err, "schema.Parse")
		return
	}
	option = kithttp.ServerBefore(
		func(ctx context.Context, r *http.Request) context.Context {
			return SetScopesToCtx(ctx, r, *tSchema)
		})
	return
}

func SetScopesToCtx(ctx context.Context, r *http.Request, tSchema schema.Schema) context.Context {
	if ctx == nil {
		ctx = r.Context()
	}
	value := CtxGormScopesValue{}
	value.QScopes, value.Err = HttpQScopes(r, tSchema)

	otherScopes := make([]func(db *gorm.DB) *gorm.DB, 0)
	sortFn, err := SortScope(r, tSchema, GetFieldByJson)
	if err != nil {
		value.Err = errors.Wrap(value.Err, err.Error())
	}
	otherScopes = append(otherScopes, sortFn)

	// pageFn, err := PagingScope(r)
	// if err != nil {
	// value.Err = errors.Wrap(value.Err, err.Error())
	// }
	// otherScopes = append(otherScopes, pageFn)

	preloadFn, err := PreloadScope(r, tSchema)
	if err != nil {
		value.Err = errors.Wrap(value.Err, err.Error())
	}
	otherScopes = append(otherScopes, preloadFn)

	selectFn, err := SelectScope(r, tSchema)
	if err != nil {
		value.Err = errors.Wrap(value.Err, err.Error())
	}
	value.OtherScopes = append(otherScopes, selectFn)
	return context.WithValue(ctx, myctx.CtxGormScopesKey, value)
}

func SetQScopes(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	value, ok := ctx.Value(myctx.CtxGormScopesKey).(CtxGormScopesValue)
	if !ok {
		return db, nil
	}
	if value.Err != nil {
		return db, value.Err
	}

	for _, fn := range value.QScopes {
		db = fn(db)
	}
	return db, nil
}

func SetScopes(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	value, ok := ctx.Value(myctx.CtxGormScopesKey).(CtxGormScopesValue)
	if !ok {
		return db, nil
	}
	if value.Err != nil {
		return db, value.Err
	}

	for _, v := range value.QScopes {
		db = v(db)
	}
	return db.Scopes(value.OtherScopes...), nil

}

func SetOtherScopes(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	value, ok := ctx.Value(myctx.CtxGormScopesKey).(CtxGormScopesValue)
	if !ok {
		return db, nil
	}
	if value.Err != nil {
		return db, value.Err
	}

	for _, fn := range value.OtherScopes {
		db = fn(db)
	}
	return db, nil
}

type GenxScopesReq struct {
	Field string
	Op    string
	Value interface{}
}

func GenxScopes(tSchema schema.Schema, req []GenxScopesReq) (fns []func(db *gorm.DB) *gorm.DB, err error) {
	for _, v := range req {
		var pq qParam
		var fn func(db *gorm.DB) *gorm.DB
		pq, err = genxParseQ(v.Field, v.Op, v.Value)
		if err != nil {
			return
		}
		fn, err = gen(pq, tSchema, GetFieldByFieldName)

		if err != nil {
			return
		}

		fns = append(fns, fn)

	}
	return
}

func HttpQScopes(r *http.Request, tSchema schema.Schema) (fns []func(db *gorm.DB) *gorm.DB, err error) {
	fns, err = httpQScope(r, tSchema, GetFieldByJson)
	if err != nil {
		err = errors.Wrap(err, "QScope")
		return
	}
	return
}

func PagingScope(r *http.Request) (fn func(db *gorm.DB) *gorm.DB, err error) {
	pageStr := r.URL.Query().Get("_page")
	pageSizeStr := r.URL.Query().Get("_pageSize")

	if pageStr == "" || pageSizeStr == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}, nil
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		err = fmt.Errorf("page参数错误: %s", pageStr)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		err = fmt.Errorf("pageSize参数错误: %s", pageSizeStr)
		return
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}, nil
}

func SelectScope(r *http.Request, tSchema schema.Schema) (fn func(db *gorm.DB) *gorm.DB, err error) {
	o := r.URL.Query().Get("_select")
	if o == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}, nil
	}

	l := strings.Split(o, ",")
	var selectList []string
	for _, field := range l {
		f, ok := tSchema.FieldsByName[field]
		if !ok {
			err = fmt.Errorf("未知的可选择字段: %s", field)
			return
		}
		dbName := f.DBName

		selectList = append(selectList, dbName)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(selectList)
	}, nil
}

func PreloadScope(r *http.Request, tSchema schema.Schema) (fn func(db *gorm.DB) *gorm.DB, err error) {
	o := r.URL.Query().Get("_include")
	if o == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}, nil
	}
	tablesList := strings.Split(o, ",")
	for _, tables := range tablesList {
		err = depthTable(tSchema, tables)
		if err != nil {
			return
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		for _, v := range tablesList {
			db = db.Preload(v)
		}
		return db
	}, nil
}

// 默认创建时间倒序
func SortScope(r *http.Request, tSchema schema.Schema, getFieldFunc func(sa *schema.Schema, name string) (*schema.Field, error)) (fn func(db *gorm.DB) *gorm.DB, err error) {
	o, ok := r.URL.Query()["_sort"]
	if !ok {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}, nil
	}

	sortList := make([]string, 0)

	for _, v := range o {
		var field string
		var order string
		if strings.HasPrefix(v, "-") {
			field = strings.TrimPrefix(v, "-")
			order = "DESC"
		} else {
			field = v
			order = "ASC"
		}

		dbField, err := getFieldFunc(&tSchema, field)
		if err != nil {
			err = errors.Wrap(err, "GetField")
			return nil, err
		}

		// f, ok := tSchema.FieldsByName[field]
		// if !ok {
		// err = fmt.Errorf("未知的可排序字段: %s", field)
		// return
		// }

		sortList = append(sortList, dbField.DBName+" "+order)
	}
	return func(db *gorm.DB) *gorm.DB {
		for _, v := range sortList {
			db = db.Order(v)
		}
		return db
	}, nil
}

func depthTable(tSchema schema.Schema, tables string) error {
	ts := strings.Split(tables, ".")
	tmpSchema := tSchema
	for _, t := range ts {
		relation, ok := tmpSchema.Relationships.Relations[t]
		if !ok {
			return errors.Errorf("未知的关联: %s/%s", tables, t)
		}
		tmpSchema = *(relation.FieldSchema)
	}
	return nil
}
