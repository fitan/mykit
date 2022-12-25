package mycrud

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

type GetOneActionMethod func(ctx context.Context, i interface{}) (interface{}, error)
type GetManyActionMethod func(ctx context.Context, i []interface{}) (interface{}, error)
type UpdateOneActionMethod func(ctx context.Context, id string, i interface{}) (interface{}, error)
type UpdateManyActionMethod func(ctx context.Context, i []interface{}) (interface{}, error)
type DeleteOneActionMethod func(ctx context.Context, id string) (interface{}, error)
type DeleteManyActionMethod func(ctx context.Context, ids string) (interface{}, error)
type CreateOneActionMethod func(ctx context.Context, i interface{}) (interface{}, error)
type CreateManyActionMethod func(ctx context.Context, i []interface{}) (interface{}, error)

type RegisterMethodHelper struct {
	crud      *CRUD
	tableName string
	methodMap map[string]methodMsg
}

func (r *RegisterMethodHelper) hasMethod(methodName string) error {
	_, ok := r.crud.tables[r.tableName].methodMap[methodName]
	if ok {
		return errors.Errorf("method %s already exists", methodName)
	}
	return nil
}

func (r *RegisterMethodHelper) Method(name string, fn interface{}, enc kithttp.EncodeResponseFunc, opts []kithttp.ServerOption) {
	if err := r.hasMethod(name); err != nil {
		panic(err)
	}

	switch t := fn.(type) {
	case GetOneActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			getOneHas: true,
			getOne:    t,
			enc:       enc,
			options:   opts,
		}
	case GetManyActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			getManyHas: true,
			getMany:    t,
			enc:        enc,
			options:    opts,
		}
	case UpdateOneActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			updateOneHas: true,
			updateOne:    t,
			enc:          enc,
			options:      opts,
		}
	case UpdateManyActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			updateManyHas: true,
			updateMany:    t,
			enc:           enc,
			options:       opts,
		}
	case DeleteOneActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			deleteOneHas: true,
			deleteOne:    t,
			enc:          enc,
			options:      opts,
		}
	case DeleteManyActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			deleteManyHas: true,
			deleteMany:    t,
			enc:           enc,
			options:       opts,
		}
	case CreateOneActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			createOneHas: true,
			createOne:    t,
			enc:          enc,
			options:      opts,
		}
	case CreateManyActionMethod:
		r.crud.tables[r.tableName].methodMap[name] = methodMsg{
			createManyHas: true,
			createMany:    t,
			enc:           enc,
			options:       opts,
		}
	default:
		panic(errors.Errorf("method %s is not a valid method", name))
	}
}
