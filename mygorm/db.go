package mygorm

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type MyKitCtxKey int

const (
	CtxDbKey MyKitCtxKey = iota + 1
)

type DB struct {
	db *gorm.DB
}

func New(db *gorm.DB) *DB {
	return &DB{db: db}
}

// ctx中包含db则返回db，如果不包含申请一个
func (d *DB) GetDb(ctx context.Context) (db *gorm.DB) {
	db, ok := ctx.Value(CtxDbKey).(*gorm.DB)
	if ok {
		return
	}

	return d.db.WithContext(ctx)
}

// 创建事务 ctx中放入db，生成一个fn，调用fn时提交或回滚事务
// 如果ctx中已经有db则使用ctx中的db,返回一个什么都不做的fn，由最开始的调用者提交或回滚事务
func (d *DB) Tx(ctx context.Context) (res context.Context, fn func(err error) (res error)) {
	_, ok := ctx.Value(CtxDbKey).(*gorm.DB)
	if ok {
		return ctx, func(err error) error {
			return err
		}
	}

	tx := d.db.Begin().WithContext(ctx)
	res = context.WithValue(res, CtxDbKey, tx)
	fn = func(err error) (res error) {
		if err == nil {
			return tx.Commit().Error
		} else {
			err1 := tx.Rollback().Error
			if err1 != nil {
				return errors.Wrap(err, err1.Error())
			}
			return err
		}
	}
	return
}

func (d *DB) Db() *gorm.DB {
	return d.db
}
