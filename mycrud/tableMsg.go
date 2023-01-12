package mycrud

import (
	"gorm.io/gorm/schema"
)

type TableMsg struct {
	oneObjFn  func() interface{}
	manyObjFn func() interface{}
	schema    schema.Schema
}
