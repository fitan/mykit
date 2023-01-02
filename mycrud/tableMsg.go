package mycrud

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"gorm.io/gorm/schema"
)

type tableMsg struct {
	oneObjFn   func() interface{}
	manyObjFn  func() interface{}
	schema     schema.Schema
	encMap     map[string]kithttp.EncodeResponseFunc
	dtoMap     map[string]func(v interface{}) interface{}
	optionsMap map[string][]kithttp.ServerOption
}

func (t *tableMsg) Dto(methodName string, dto func(v interface{}) interface{}) *tableMsg {
	t.dtoMap[methodName] = dto
	return t
}

func (t *tableMsg) Option(methodName string, option kithttp.ServerOption) *tableMsg {
	t.optionsMap[methodName] = append(t.optionsMap[methodName], option)
	return t
}

func (t *tableMsg) Encode(methodName string, encode kithttp.EncodeResponseFunc) *tableMsg {
	t.encMap[methodName] = encode
	return t
}
