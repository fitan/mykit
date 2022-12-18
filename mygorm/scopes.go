package mygorm

import "gorm.io/gorm"

func WhereEqScope(i interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if i != nil {
			return db.Where(i)
		}
		return db
	}
}

// 分页对象
type Paging struct {
	// 页码
	Page int `json:"page"`
	// 每页数量
	PageSize int `json:"pageSize"`
}

func PagingScope(p Paging) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.PageSize != 0 && p.Page != 0 {
			offset := (p.Page - 1) * p.PageSize
			return db.Offset(offset).Limit(p.PageSize)
		}

		return db
	}
}

type Order struct {
	Order string
}

// 默认创建时间倒序
func OrderScope(o Order) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if o.Order == "" {
			return db.Order("created_at desc")
		}
		return db
	}
}
