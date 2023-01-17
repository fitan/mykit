package mygorm

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"regexp"
	"strings"
)

var ops = map[string]string{
	"=":   "= ?",
	"!=":  "!= ?",
	">":   "> ?",
	"<":   "< ?",
	">=":  ">= ?",
	"<=":  "<= ?",
	"?=":  "in ?",
	"!?=": "not in ?",
	"~=":  "like ?",
	"!~=": "not like ?",
	"><":  "between ? and ?",
	"!><": "not between ? and ?",
	//"isnull": "=null",
	//"notnull": "!=null",
}

func QScope(r *http.Request, tSchema schema.Schema) (fns []func(db *gorm.DB) *gorm.DB, err error) {
	qList, ok := r.URL.Query()["_q"]
	if !ok {
		return
	}
	reg, _ := regexp.Compile(`[!?~=><]+`)

	for _, v := range qList {
		var pq qParam
		var fn func(db *gorm.DB) *gorm.DB
		op := reg.FindString(v)
		pq, err = parseQ(v, op)
		if err != nil {
			return
		}

		fn, err = gen(pq, tSchema)
		if err != nil {
			return
		}

		fns = append(fns, fn)
	}

	return
}

type qParam struct {
	field string
	op    string
	value string
	sqlOp string
}

func (s qParam) toScope(fieldName string) (res func(db *gorm.DB) *gorm.DB, err error) {
	switch s.op {
	case "=", "!=", ">", "<", ">=", "<=", "~=", "!~=":
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fieldName+" "+s.sqlOp, s.value)
		}, nil
	case "?=", "!?=":
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fieldName+" "+s.sqlOp, strings.Split(s.value, ","))
		}, nil
	case "><", "<>":
		l := strings.SplitN(s.value, ",", 2)
		if len(l) != 2 {
			err = fmt.Errorf("wrong format %s", s.value)
			return
		}
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fieldName+" "+s.sqlOp, l[0], l[1])
		}, nil
	default:
		return nil, fmt.Errorf("toSqlValue not found op: %s", s.op)
	}
}

func (s qParam) localTable() bool {
	return !strings.Contains(s.field, ".")
}

type relationTable struct {
	tableName         string
	relationTableName string
	foreignKey        string
	primaryKey        string
}

func gen(param qParam, tSchema schema.Schema) (fn func(db *gorm.DB) *gorm.DB, err error) {
	var relationTables []relationTable
	var tmpSchema *schema.Schema
	var tmpField *schema.Field
	var ok bool

	tmpSchema = &tSchema

	// table1.table2.field1
	fieldList := strings.Split(param.field, ".")
	for i, v := range fieldList {
		if len(fieldList)-1 == i {
			tmpField, ok = tmpSchema.FieldsByName[v]
			if !ok {
				err = fmt.Errorf("not found field: %s", v)
			}
			break
		}
		if tmpSchema == nil {
			err = fmt.Errorf("field %s schema is null", v)
			return
		}

		relation, ok := tmpSchema.Relationships.Relations[v]
		if !ok {
			err = fmt.Errorf("not found relation: %s", v)
			return
		}

		if len(relation.References) == 0 {
			err = fmt.Errorf("not found reference: %s", v)
			return
		}

		relationTables = append(relationTables, relationTable{
			tableName:         tmpSchema.Table,
			relationTableName: relation.FieldSchema.Table,
			foreignKey:        relation.References[0].ForeignKey.DBName,
			primaryKey:        relation.References[0].PrimaryKey.DBName,
		})

		//fmt.Println("tmpField", tmpSchema.FieldsByName, v)
		//fmt.Println("tmpfieldValue", tmpSchema.FieldsByName[v])
		//spew.Dump(tmpSchema.FieldsByName[v].Schema)
		//fmt.Println("tmpfieldValueschema", tmpSchema.FieldsByName[v].Schema)
		tmpField, ok = tmpSchema.FieldsByName[v]
		if !ok {
			err = fmt.Errorf("not found field: %s", v)
			return
		}
		//fmt.Println("tmpfieldvalue", tmpField)
		//spew.Dump(tmpField)
		//if !ok {
		//	err = fmt.Errorf("not found field: %s", v)
		//	return
		//}]
		tmpSchema = relation.FieldSchema
	}

	fn, err = param.toScope(tmpField.DBName)
	if err != nil {
		return
	}

	spew.Dump(relationTables)

	for i := len(relationTables) - 1; i >= 0; i-- {
		r := relationTables[i]

		tmpFn := fn

		fn = func(db *gorm.DB) *gorm.DB {
			value := tmpFn(db.Session(&gorm.Session{NewDB: true}).Table(r.relationTableName)).Select(r.foreignKey)
			return db.Where(r.primaryKey+" IN (?)", value)
		}

	}

	return
}

func parseQ(s, op string) (res qParam, err error) {
	sqlOp, ok := ops[op]
	if !ok {
		err = fmt.Errorf("not found op: %s", op)
		return
	}
	l := strings.SplitN(s, op, 2)
	if len(l) != 2 {
		err = fmt.Errorf("wrong format %s", s)
		return
	}

	res.op = op
	res.field = l[0]
	res.value = l[1]
	res.sqlOp = sqlOp
	return
}
