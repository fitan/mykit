package mycrud

import (
	"fmt"
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
		want *Core
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
				crud := NewCRUD(tt.args.m, tt.args.db, tt.args.encode, nil)
				userReg, err := crud.RegisterTable(
					func() interface{} {
						return &ormdata.User{}
					}, func() interface{} {
						return &[]ormdata.User{}
					})
				if err != nil {
					panic(err)
				}
				_, err = crud.RegisterTable(
					func() interface{} {
						return &ormdata.Todo{}
					}, func() interface{} {
						return &[]ormdata.Todo{}
					})
				if err != nil {
					panic(err)
				}
				postReg, err := crud.RegisterTable(
					func() interface{} {
						return &ormdata.Post{}
					},
					func() interface{} {
						return &[]ormdata.Post{}
					},
				)
				if err != nil {
					panic(err)
				}
				_, err = crud.RegisterTable(
					func() interface{} {
						return &ormdata.Photo{}
					},
					func() interface{} {
						return &[]ormdata.Photo{}
					},
				)
				if err != nil {
					panic(err)
				}

				_, err = crud.RegisterTable(
					func() interface{} {
						return &ormdata.Comment{}
					},
					func() interface{} {
						return &[]ormdata.Comment{}
					},
				)
				if err != nil {
					panic(err)
				}
				_, err = crud.RegisterTable(
					func() interface{} {
						return &ormdata.Album{}
					},
					func() interface{} {
						return &[]ormdata.Album{}
					},
				)
				if err != nil {
					panic(err)
				}
				fmt.Println(userReg, postReg)

				//userReg.Dto(GetOneMethodName, func(v interface{}) interface{} {
				//	data := v.(*ormdata.User)
				//	return UserDto{ID:  data.ID}
				//}).Dto(GetManyMethodName, func(v interface{}) interface{} {
				//	data := v.(GetManyData)
				//	var result []UserDto
				//	for _, v := range *(data.List.(*[]ormdata.User)) {
				//		result = append(result, UserDto{ID: v.ID})
				//	}
				//
				//	data.List = result
				//	return data
				//})
				//
				//postReg.Dto(GetOneMethodName, func(v interface{}) interface{} {
				//	data := v.(*ormdata.Post)
				//	return UserDto{ID: data.ID}
				//}).Dto(GetManyMethodName, func(v interface{}) interface{} {
				//	data := v.(GetManyData)
				//	var result []UserDto
				//	for _, v := range *(data.List.(*[]ormdata.Post)) {
				//		result = append(result, UserDto{ID: v.ID})
				//	}
				//
				//	data.List = result
				//	return data
				//})

				crud.run()
				crud.D2Handler(crud.m)
				log.Println("crud run")
				r.Run(":8080")
			},
		)
	}
}
