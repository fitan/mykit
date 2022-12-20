package mygorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"strconv"
	"strings"
)

func PagingScope(r *http.Request) (fn func(db *gorm.DB) *gorm.DB, err error) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

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

// 默认创建时间倒序
func SortScope(r *http.Request, tSchema schema.Schema) (fn func(db *gorm.DB) *gorm.DB, err error) {
	o, ok := r.URL.Query()["sort"]
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
