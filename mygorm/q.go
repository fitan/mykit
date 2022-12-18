package mygorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// [!?~=><]+
//var ops = map[string]string{
//	"eq": "=",
//	"ne": "!=",
//	"gt": ">",
//	"lt": "<",
//	"ge": ">=",
//	"le": "<=",
//	"in": "?=",
//	"nin": "!?=",
//	"like": "~=",
//	"nlike": "!~=",
//	"isnull": "=null",
//	"notnull": "!=null",
//}

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
	//"isnull": "=null",
	//"notnull": "!=null",
}

func Q(r *http.Request, t interface{}) (fns []func(db *gorm.DB) *gorm.DB, err error) {
	tSchema, err := schema.Parse(t, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return
	}
	qList, ok := r.URL.Query()["q"]
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

		fn, err = gen(pq, *tSchema)
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

func (s qParam) toSqlValue() (res interface{}, err error) {
	switch s.op {
	case "=", "!=", ">", "<", ">=", "<=", "~=", "!~=":
		return s.value, nil
	case "?=", "!?=":
		return strings.Split(s.value, "."), nil
	default:
		return nil, fmt.Errorf("toSqlValue not found op: %s", s.op)
	}
}

func (s qParam) localTable() bool {
	return !strings.Contains(s.field, ".")
}

type relationTable struct {
	tableName  string
	foreignKey string
	primaryKey string
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
		relationTables = append(relationTables, relationTable{
			tableName:  tmpSchema.FieldsByName[v].DBName,
			foreignKey: tmpSchema.Relationships.Relations[v].References[0].ForeignKey.DBName,
			primaryKey: tmpSchema.Relationships.Relations[v].References[0].PrimaryKey.DBName,
		})

		//fmt.Println("tmpField", tmpSchema.FieldsByName, v)
		//fmt.Println("tmpfieldValue", tmpSchema.FieldsByName[v])
		//spew.Dump(tmpSchema.FieldsByName[v].Schema)
		//fmt.Println("tmpfieldValueschema", tmpSchema.FieldsByName[v].Schema)
		tmpField = tmpSchema.FieldsByName[v]
		//fmt.Println("tmpfieldvalue", tmpField)
		//spew.Dump(tmpField)
		//if !ok {
		//	err = fmt.Errorf("not found field: %s", v)
		//	return
		//}]
		tmpSchema = tmpSchema.Relationships.Relations[v].FieldSchema
	}

	tableName := tmpSchema.Table

	sqlValue, err := param.toSqlValue()
	if err != nil {
		return
	}

	fn = func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{NewDB: true}).Table(tableName).Where(tmpField.DBName+" "+param.sqlOp, sqlValue)
	}

	for i := len(relationTables) - 1; i >= 0; i-- {
		r := relationTables[i]
		fmt.Println(r)

		tmpFn := func(db *gorm.DB) *gorm.DB {
			return db.Session(&gorm.Session{NewDB: true}).Table(r.tableName).Where(r.foreignKey+" = ?", fn(db).Select(r.primaryKey))
		}

		fn = tmpFn
	}

	return
}

func parseQ(s, op string) (res qParam, err error) {
	sqlOp, ok := ops[op]
	if !ok {
		err = fmt.Errorf("not found op: %s", op)
		return
	}
	l := strings.Split(s, op)
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
