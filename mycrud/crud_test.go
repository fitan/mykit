package mycrud

import (
	"github.com/fitan/mykit/mycrud/ormdata"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestNewCRUD(t *testing.T) {
	db := ormdata.New().Debug()
	m := mux.NewRouter()
	type args struct {
		m      *mux.Router
		db     *gorm.DB
		encode kithttp.EncodeResponseFunc
	}
	tests := []struct {
		name string
		args args
		want *CRUD
	}{
		{
			name: "test",
			args: args{
				m:      m,
				db:     db,
				encode: nil,
			},
			want: nil,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				crud := NewCRUD(tt.args.m, tt.args.db, tt.args.encode)
				crud.RegisterTable(
					func() interface{} {
						return ormdata.User{}
					}, func() interface{} {
						return make([]ormdata.User, 0)
					})
				crud.RegisterTable(
					func() interface{} {
						return ormdata.Todo{}
					}, func() interface{} {
						return make([]ormdata.Todo, 0)
					})
				crud.RegisterTable(
					func() interface{} {
						return ormdata.Post{}
					},
					func() interface{} {
						return make([]ormdata.Post, 0)
					},
				)
				crud.RegisterTable(
					func() interface{} {
						return ormdata.Photo{}
					},
					func() interface{} {
						return make([]ormdata.Photo, 0)
					},
				)
				crud.RegisterTable(
					func() interface{} {
						return ormdata.Comment{}
					},
					func() interface{} {
						return make([]ormdata.Comment, 0)
					},
				)
				crud.RegisterTable(
					func() interface{} {
						return ormdata.Album{}
					},
					func() interface{} {
						return make([]ormdata.Album, 0)
					},
				)

				crud.run()
				http.ListenAndServe(":8080", crud.m)
			},
		)
	}
}
