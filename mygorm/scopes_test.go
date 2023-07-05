package mygorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"sync"
	"testing"
)

func TestGenxScopes(t *testing.T) {

	dsn := "spider_dev:spider_dev123@tcp(10.170.34.22:3307)/spider_dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	s, _ := schema.Parse(&PhysicalMachine{}, &sync.Map{}, schema.NamingStrategy{})

	type args struct {
		tSchema schema.Schema
		req     []GenxScopesReq
	}
	tests := []struct {
		name    string
		args    args
		wantFns []func(db *gorm.DB) *gorm.DB
		wantErr bool
	}{
		{
			name: "genx",
			args: args{
				tSchema: *s,
				req: []GenxScopesReq{
					{
						Field: "Brand.ProductType",
						Op:    "=",
						Value: "服务器",
					},
					{
						Field: "UUID",
						Op:    "=",
						Value: "83c63f28970d433597f6caf2696ceab4",
					},
					{
						Field: "Brand.Users.Name",
						Op:    "=",
						Value: "张三",
					},
					{
						Field: "Brand.ID",
						Op:    ">",
						Value: "10",
					},
					{
						Field: "Brand.UUID",
						Op:    "?=",
						Value: []string{"83c63f28970d433597f6caf2696ceab4", "83c63f28970d433597f6caf2696ceab5"},
					},
					{
						Field: "Brand.CreatedAt",
						Op:    "><",
						Value: []string{"2021-01-01", "2021-01-02"},
					},
				},
			},
			wantFns: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFns, err := GenxScopes(tt.args.tSchema, tt.args.req)
			if err != nil {
				panic(err)
			}
			data := make([]PhysicalMachine, 0)
			tmpDB := db.Session(&gorm.Session{DryRun: true})
			fmt.Println(tmpDB.Scopes(gotFns...).Find(&data).Statement.SQL.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("GenxScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFns, tt.wantFns) {
				t.Errorf("GenxScopes() gotFns = %v, want %v", gotFns, tt.wantFns)
			}
		})
	}
}
