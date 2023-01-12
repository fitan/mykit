package mycrud

import (
	"context"
	"fmt"
	"github.com/fitan/mykit/mygorm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
)

type RepoImpl interface {
	GetOne(ctx context.Context, id string) (data interface{}, err error)
	GetMany(ctx context.Context, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error)

	CreateOne(ctx context.Context, body interface{}) (err error)
	CreateMany(ctx context.Context, data interface{}) (err error)

	UpdateOne(ctx context.Context, id string, data interface{}) (err error)
	UpdateMany(ctx context.Context, data interface{}) (err error)

	DeleteOne(ctx context.Context, id string) (err error)
	DeleteMany(ctx context.Context, ids []string) (err error)

	GetRelationOne(ctx context.Context, id string, scopes []func(db *gorm.DB) *gorm.DB) (data interface{}, err error)
	GetRelationMany(ctx context.Context, id string, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error)

	CreateRelationOne(ctx context.Context, id string, body interface{}) (err error)
	CreateRelationMany(ctx context.Context, id string, body interface{}) (err error)
}

type Repo struct {
	Core     *Core
	TableMsg *TableMsg
	//RelationTableMsg *GetTableMsg
}

func (r *Repo) GetOne(ctx context.Context, id string) (data interface{}, err error) {
	db := r.Core.db.Db(ctx)

	obj := r.TableMsg.oneObjFn()
	err = db.Model(r.TableMsg.oneObjFn()).Where("id = ?", id).First(obj).Error
	return obj, err
}

func (r *Repo) GetMany(ctx context.Context, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
	var total int64
	totalDB := r.Core.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Model(r.TableMsg.oneObjFn()).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Count")
		return
	}
	db := r.Core.db.Db(ctx).Model(r.TableMsg.oneObjFn()).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}
	list := r.TableMsg.manyObjFn()
	err = db.Find(list).Error
	if err != nil {
		err = errors.Wrap(err, "db.Find")
		return
	}

	data.Total = total
	data.List = list
	return
}

func (r *Repo) CreateOne(ctx context.Context, data interface{}) (err error) {
	db, commit := r.Core.db.Tx(ctx)
	defer commit(err)

	err = db.Model(r.TableMsg.oneObjFn()).Create(data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Create")
		return
	}
	return
}

func (r *Repo) CreateMany(ctx context.Context, data interface{}) (err error) {
	db, commit := r.Core.db.Tx(ctx)
	defer commit(err)

	err = db.Model(r.TableMsg.oneObjFn()).CreateInBatches(data, 20).Error
	if err != nil {
		err = errors.Wrap(err, "db.CreateInBatches")
		return
	}
	return
}

func (r *Repo) UpdateOne(ctx context.Context, id string, data interface{}) (err error) {
	db, commit := r.Core.db.Tx(ctx)
	defer commit(err)

	err = db.Model(r.TableMsg.oneObjFn()).Select("*").Where("id = ?", id).Updates(data).Error
	if err != nil {
		err = errors.Wrap(err, "db.Updates()")
		return
	}
	return
}

func (r *Repo) UpdateMany(ctx context.Context, data interface{}) (err error) {
	db, commit := r.Core.db.Tx(ctx)
	defer commit(err)

	refV := reflect.ValueOf(data)

	if refV.Kind() == reflect.Ptr {
		refV = refV.Elem()
	}

	switch refV.Kind() {
	case reflect.Slice:
		for i := 0; i < refV.Len(); i++ {
			err = db.Model(refV.Index(i).Interface()).Updates(refV.Index(i).Interface()).Error
			if err != nil {
				err = errors.Wrap(err, "db.Updates()")
				return
			}
		}
	default:
		return errors.New("data must be slice")
	}

	return
}

func (r *Repo) DeleteOne(ctx context.Context, id string) (err error) {
	db, commit := r.Core.db.Tx(ctx)
	defer commit(err)

	err = db.Model(r.TableMsg.oneObjFn()).Where("id = ?", id).Delete(r.TableMsg.oneObjFn()).Error
	return
}

func (r *Repo) DeleteMany(ctx context.Context, ids []string) (err error) {
	db, commit := r.Core.db.Tx(ctx)
	defer commit(err)

	err = db.Model(r.TableMsg.oneObjFn()).Where("id in (?)", ids).Delete(r.TableMsg.oneObjFn()).Error
	return
}

func (r *Repo) GetRelationOne(
	ctx context.Context,
	id string,
	relationTableName string,
	scopes []func(db *gorm.DB) *gorm.DB,
) (data interface{}, err error) {
	relationTableMsg, err := r.Core.GetTableMsg(relationTableName)
	if err != nil {
		return
	}

	var relation schema.Relationship
	var hasRelation bool

	for _, v := range relationTableMsg.schema.Relationships.Relations {
		if v.FieldSchema.Table == r.TableMsg.schema.Table {
			relation = *v
			//_, err = r.Core.GetTableMsg(relationTableName)
			//if err != nil {
			//	return
			//}
			hasRelation = true
			break
		}
	}

	if hasRelation == false {
		err = errors.New("no relation")
		return
	}

	if len(relation.References) == 0 {
		err = fmt.Errorf("not found reference: %s", r.TableMsg.schema.Table)
		return
	}

	relationForeignKey := relation.References[0].ForeignKey.DBName
	relationPrimaryKey := relation.References[0].PrimaryKey.DBName

	data = relationTableMsg.oneObjFn()

	db := r.Core.db.Db(ctx)
	db.Model(relationTableMsg.oneObjFn()).Where(relationPrimaryKey+" in (?)", db.Session(&gorm.Session{NewDB: true}).Table(r.TableMsg.schema.Table).Select(relationForeignKey).Where(relationForeignKey+" = ?", id)).First(data)

	return data, nil
}

func (r *Repo) GetRelationMany(ctx context.Context, id string, relationTableName string, scopes []func(db *gorm.DB) *gorm.DB) (data GetManyData, err error) {
	relationTableMsg, err := r.Core.GetTableMsg(relationTableName)
	if err != nil {
		return
	}

	var relation schema.Relationship
	var hasRelation bool

	for _, v := range relationTableMsg.schema.Relationships.Relations {
		if v.FieldSchema.Table == r.TableMsg.schema.Table {
			relation = *v
			//_, err = r.Core.GetTableMsg(relationTableName)
			//if err != nil {
			//	return
			//}
			hasRelation = true
			break
		}
	}

	if hasRelation == false {
		err = errors.New("no relation")
		return
	}

	if len(relation.References) == 0 {
		err = fmt.Errorf("not found reference: %s", r.TableMsg.schema.Table)
		return
	}

	relationForeignKey := relation.References[0].ForeignKey.DBName
	relationPrimaryKey := relation.References[0].PrimaryKey.DBName

	var total int64
	totalDB := r.Core.db.Db(ctx)
	totalDB, err = mygorm.SetQScopes(ctx, totalDB)
	if err != nil {
		return
	}
	err = totalDB.Model(relationTableMsg.oneObjFn()).Where(relationForeignKey+" in (?)", totalDB.Session(&gorm.Session{NewDB: true}).Table(r.TableMsg.schema.Table).Select(relationPrimaryKey).Where(relationPrimaryKey+" = ?", id)).Scopes(scopes...).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "db.Count")
		return
	}

	db := r.Core.db.Db(ctx).Model(relationTableMsg.oneObjFn()).Scopes(scopes...)
	db, err = mygorm.SetScopes(ctx, db)
	if err != nil {
		return
	}

	list := relationTableMsg.manyObjFn()

	db.Where(relationForeignKey+" in (?)", db.Session(&gorm.Session{NewDB: true}).Table(r.TableMsg.schema.Table).Select(relationPrimaryKey).Where(relationPrimaryKey+" = ?", id)).Find(list)

	data.List = list
	data.Total = total

	return data, nil
}

func (r *Repo) CreateRelationOne(
	ctx context.Context,
	id, relationTableName string,
	body interface{},
) (err error) {
	return r.createRelation(ctx, id, relationTableName, body)
}

func (r *Repo) CreateRelationMany(
	ctx context.Context,
	id, relationTableName string,
	body interface{},
) (err error) {
	return r.createRelation(ctx, id, relationTableName, body)
}

func (r *Repo) createRelation(ctx context.Context, id string, relationTableName string, body interface{}) (err error) {
	relationTableMsg, err := r.Core.GetTableMsg(relationTableName)
	if err != nil {
		return
	}

	db := r.Core.db.Db(ctx)

	var hasRelation bool
	var relationFieldName string

	for k, v := range r.TableMsg.schema.Relationships.Relations {
		if v.FieldSchema.Table == relationTableMsg.schema.Table {
			hasRelation = true
			relationFieldName = k
		}
	}

	if !hasRelation {
		err = fmt.Errorf("table %s has no relation with table %s", r.TableMsg.schema.Table, relationTableMsg.schema.Table)
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
	switch r.TableMsg.schema.FieldsByDBName["id"].FieldType.Kind().String() {
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
		err = fmt.Errorf("not support id type %s", r.TableMsg.schema.FieldsByDBName["id"].DataType)
		return
	}

	model := r.TableMsg.oneObjFn()
	reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(gormID))

	err = db.Model(model).Association(relationFieldName).Append(body)
	if err != nil {
		err = errors.Wrap(err, "db.Append")
		return
	}

	return
}
