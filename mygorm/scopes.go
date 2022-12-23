package mygorm

import (
	"context"
	"fmt"
	"github.com/fitan/mykit/myctx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func SetScopes(ctx context.Context, db *gorm.DB) *gorm.DB {
	scopes, ok := ctx.Value(myctx.CtxGormScopesKey).([]func(db *gorm.DB) *gorm.DB)
	if !ok {
		return db
	}
	for _, fn := range scopes {
		db = fn(db)
	}
	return db
}

func Scopes(r *http.Request, i interface{}) (r2 *http.Request, err error) {
	tSchema, err := schema.Parse(i, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return
	}

	fns, err := QScope(r, *tSchema)
	if err != nil {
		err = errors.Wrap(err, "QScope")
		return
	}

	sortFn, err := SortScope(r, *tSchema)
	if err != nil {
		err = errors.Wrap(err, "SortScope")
		return
	}
	fns = append(fns, sortFn)

	pageFn, err := PagingScope(r)
	if err != nil {
		err = errors.Wrap(err, "PagingScope")
		return
	}
	fns = append(fns, pageFn)

	selectFn, err := SelectScope(r, *tSchema)
	if err != nil {
		err = errors.Wrap(err, "SelectScope")
		return
	}

	fns = append(fns, selectFn)

	ctx := context.WithValue(r.Context(), myctx.CtxGormScopesKey, fns)
	return r.WithContext(ctx), nil
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

// 默认创建时间倒序
func SortScope(r *http.Request, tSchema schema.Schema) (fn func(db *gorm.DB) *gorm.DB, err error) {
	o, ok := r.URL.Query()["_sort"]
	if !ok {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}, nil
	}

	sortList := make([]string, 0)

	for _, v := range o {
		l := strings.SplitN(v, ",", 2)
		if len(l) != 2 {
			err = fmt.Errorf("sort参数错误: %s", v)
			return
		}
		field := l[0]
		order := l[1]
		if order != "asc" && order != "desc" {
			err = fmt.Errorf("sort参数错误: %s", v)
			return
		}
		f, ok := tSchema.FieldsByName[field]
		if !ok {
			err = fmt.Errorf("未知的可排序字段: %s", field)
			return
		}
		dbName := f.DBName

		sortList = append(sortList, dbName+" "+order)
	}
	return func(db *gorm.DB) *gorm.DB {
		for _, v := range sortList {
			db = db.Order(v)
		}
		return db
	}, nil
}
