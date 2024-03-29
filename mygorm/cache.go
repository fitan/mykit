package mygorm

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"sync"
)

type Cache struct {
	SchemaByType       map[reflect.Type]*schema.Schema
	FieldByJson        map[reflect.Type]map[string]*schema.Field
	TypeModelInterface map[reflect.Type]interface{}
	m                  *sync.Mutex
}

var cache *Cache

func GetSchema(i any) (*schema.Schema, error) {
	return cache.schema(i)
}

type GetFieldFunc func(sa *schema.Schema, name string) (*schema.Field, error)

func GetFieldByFieldName(sa *schema.Schema, name string) (*schema.Field, error) {
	if f, ok := sa.FieldsByName[name]; ok {
		return f, nil
	} else {
		return nil, fmt.Errorf("field %s not found", name)
	}
}

func GetFieldByJson(sa *schema.Schema, name string) (*schema.Field, error) {
	return cache.field(sa, name)
}

func GetModelInterface(p reflect.Type) interface{} {
	return cache.getModelInterface(p)
}

func init() {
	cache = &Cache{
		SchemaByType:       map[reflect.Type]*schema.Schema{},
		FieldByJson:        map[reflect.Type]map[string]*schema.Field{},
		TypeModelInterface: map[reflect.Type]interface{}{},
		m:                  &sync.Mutex{},
	}
}

func (c *Cache) getModelInterface(p reflect.Type) interface{} {
	c.m.Lock()
	defer c.m.Unlock()
	if model, ok := c.TypeModelInterface[p]; ok {
		return model
	} else {
		i := reflect.New(p).Interface()
		c.TypeModelInterface[p] = i
		return i
	}
}

func (c *Cache) schema(i any) (*schema.Schema, error) {
	c.m.Lock()
	defer c.m.Unlock()
	if tSchema, ok := c.SchemaByType[reflect.TypeOf(i)]; ok {
		return tSchema, nil
	}
	tSchema, err := schema.Parse(i, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return tSchema, errors.Wrap(err, "schema.Parse")
	}

	c.SchemaByType[tSchema.ModelType] = tSchema
	return tSchema, nil
}

func (c *Cache) field(sa *schema.Schema, name string) (*schema.Field, error) {
	c.m.Lock()
	defer c.m.Unlock()
	m, ok := c.FieldByJson[sa.ModelType]
	if !ok {
		c.FieldByJson[sa.ModelType] = make(map[string]*schema.Field)
	} else {
		field, ok := m[name]
		if ok {
			return field, nil
		}
	}

	s, ok := c.SchemaByType[sa.ModelType]
	if !ok {
		c.SchemaByType[sa.ModelType] = sa
		s = sa
	}

	return c.fieldByJson(s, name)
}

func (c *Cache) fieldByJson(sa *schema.Schema, name string) (*schema.Field, error) {
	for _, f := range sa.FieldsByName {
		j, ok := f.Tag.Lookup("json")
		if ok && j != "" {
			jsonName := strings.Split(j, ",")[0]
			if name == jsonName {
				c.FieldByJson[sa.ModelType][name] = f
				return f, nil
			}
		}
	}

	f, ok := sa.FieldsByName[name]
	if !ok {
		return nil, errors.New("field not found")
	}

	c.FieldByJson[sa.ModelType][name] = f
	return f, nil
}
