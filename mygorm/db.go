package mygorm

import (
	"context"
	"github.com/fitan/mykit/myctx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func New(db *gorm.DB) *DB {
	return &DB{db: db}
}

// ctx中包含db则返回db，如果不包含申请一个
func (d *DB) Db(ctx context.Context) (db *gorm.DB) {
	db, ok := ctx.Value(myctx.CtxGormDbKey).(*gorm.DB)
	if ok {
		return
	}

	return d.db.WithContext(ctx)
}

func (d *DB) Tx(ctx context.Context) (db *gorm.DB, commit func(err error) (res error)) {
	tmpDb, ok := ctx.Value(myctx.CtxGormDbKey).(*gorm.DB)
	if ok {
		return tmpDb, func(err error) (res error) {
			return err
		}
	}

	tx := d.db.Begin().WithContext(ctx)
	return tx, func(err error) (res error) {
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
}

// 创建事务 ctx中放入db，生成一个fn，调用fn时提交或回滚事务
// 如果ctx中已经有db则使用ctx中的db,返回一个什么都不做的fn，由最开始的调用者提交或回滚事务
func (d *DB) SetTx2Ctx(ctx context.Context) (res context.Context, commit func(err error) (res error)) {
	_, ok := ctx.Value(myctx.CtxGormDbKey).(*gorm.DB)
	if ok {
		return ctx, func(err error) error {
			return err
		}
	}

	tx := d.db.Begin().WithContext(ctx)
	res = context.WithValue(res, myctx.CtxGormDbKey, tx)
	commit = func(err error) (res error) {
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
