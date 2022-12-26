package mycrud

import (
	"github.com/fitan/mykit/mycrud/ormdata"
	"github.com/fitan/mykit/myrouter"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestNewCRUD(t *testing.T) {
	m := mux.NewRouter()
	r := myrouter.New(m)
	db := ormdata.New().Debug()
	//user := ormdata.User{}
	//err := db.Preload("Posts").Preload("Albums").Preload("Todos").Where("id = ?", 1).Find(&user).Error
	//if err != nil {
	//	panic(err)
	//}
	//json.NewEncoder(os.Stdout).Encode(user)
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
						return &ormdata.User{}
					}, func() interface{} {
						return &[]ormdata.User{}
					})
				crud.RegisterTable(
					func() interface{} {
						return &ormdata.Todo{}
					}, func() interface{} {
						return &[]ormdata.Todo{}
					})
				crud.RegisterTable(
					func() interface{} {
						return &ormdata.Post{}
					},
					func() interface{} {
						return &[]ormdata.Post{}
					},
				)
				crud.RegisterTable(
					func() interface{} {
						return &ormdata.Photo{}
					},
					func() interface{} {
						return &[]ormdata.Photo{}
					},
				)
				crud.RegisterTable(
					func() interface{} {
						return &ormdata.Comment{}
					},
					func() interface{} {
						return &[]ormdata.Comment{}
					},
				)
				crud.RegisterTable(
					func() interface{} {
						return &ormdata.Album{}
					},
					func() interface{} {
						return &[]ormdata.Album{}
					},
				)

				crud.run()
				log.Println("crud run")
				r.Run(":8080")
			},
		)
	}
}
