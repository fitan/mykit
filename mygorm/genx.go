package mygorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

type GenxScopesReq struct {
	Field string
	Op    string
	Value interface{}
}

func GenxScopes(i any, req []GenxScopesReq) (fns []func(db *gorm.DB) *gorm.DB, err error) {
	for _, v := range req {
		var pq qParam
		var fn func(db *gorm.DB) *gorm.DB
		var tSchema *schema.Schema
		pq, err = genxParseQ(v.Field, v.Op, v.Value)
		if err != nil {
			return
		}
		tSchema, err = GetSchema(i)
		if err != nil {
			return
		}
		fn, err = gen(pq, *tSchema, GetFieldByFieldName)

		if err != nil {
			return
		}

		fns = append(fns, fn)

	}
	return
}

func genxParseQ(path string, op string, value interface{}) (res qParam, err error) {
	sqlOp, ok := ops[op]
	if !ok {
		err = fmt.Errorf("not found op: %s", op)
		return
	}

	res.op = op
	res.field = path
	res.sqlOp = sqlOp

	switch res.op {
	case "=", "!=", ">", "<", ">=", "<=", "~=", "!~=":
		res.value = append(res.value, value)
	case "?=", "!?=", "><", "<>":
		vt := reflect.ValueOf(value)
		if vt.Type().Kind() == reflect.Ptr {
			vt = vt.Elem()
		}

		if vt.Kind() != reflect.Slice {
			err = fmt.Errorf("wrong format %s", value)
			return
		}

		for i := 0; i < vt.Len(); i++ {
			res.value = append(res.value, vt.Index(i).Interface())
		}
	default:
		err = fmt.Errorf("not found op: %s", op)
	}
	return
}
