package mygorm

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PhysicalMachine struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	//品牌信息
	//Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`
	//型号
	//BrandModel string `gorm:"column:brand_model;notnull;comment:'型号'" json:"brandModel"`

	// 是否采集完成

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`
}

func (PhysicalMachine) TableName() string {
	return "assets_physical_machine"
}

type Brand struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 品牌
	Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`

	// 产品类型 1:服务器 2:交换机 3:路由器 4:防火墙 5:负载均衡器 6:存储设备 7: cpu 8:内存 9:硬盘 10:网卡 10: 系统
	ProductType string `gorm:"column:product_type;notnull;comment:'产品类型'" json:"productType"`

	// 产品型号
	ProductModel string `gorm:"column:product_model;notnull;comment:'产品型号'" json:"productModel"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	Users []User `gorm:"foreignKey:UUID;references:UUID" json:"users"`
}

type User struct {
	gorm.Model
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	Name string `gorm:"column:name;notnull;comment:'姓名'" json:"name"`
}

func (b *Brand) TableName() string {
	return "assets_brand"
}

func TestScopes(t *testing.T) {
	// dsn := "root:123456@tcp(172.29.107.199:3306)/gteml?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	// dsn := "spider_dev:spider_dev123@tcp(10.170.34.22:3307)/spider_dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	s, _ := schema.Parse(&PhysicalMachine{}, &sync.Map{}, schema.NamingStrategy{})

	r := NewQuery().
		Q("brand.productType", "=", "服务器").
		Q("uuid", "=", "83c63f28970d433597f6caf2696ceab4").
		Q("brand.users.name", "=", "张三").
		Q("brand.id", ">", "10").
		Q("brand.uuid", "?=", "83c63f28970d433597f6caf2696ceab4", "83c63f28970d433597f6caf2696ceab5").
		Q("brand.uuid", "!?=", "83c63f28970d433597f6caf2696ceab4", "83c63f28970d433597f6caf2696ceab5").
		Q("brand.createdAt", "><", "2021-01-01", "2021-01-02").
		Sort("+ID", "-UUID").
		Paging("1", "10").
		NewRequest()

	//v := url.Values{}
	//v.Add("_q", "Brand.ProductType=服务器")
	//v.Add("_q", "UUID=83c63f28970d433597f6caf2696ceab4")
	//v.Add("_q", "Brand.Users.Name=张三")
	//v.Add("_q", "Brand.ID>10")
	//v.Add("_q", "Brand.UUID?=83c63f28970d433597f6caf2696ceab4,83c63f28970d433597f6caf2696ceab5")
	//v.Add("_q", "Brand.UUID!?=83c63f28970d433597f6caf2696ceab4,83c63f28970d433597f6caf2696ceab5")
	//v.Add("_q", "Brand.CreatedAt<>2021-01-01,2021-01-02")
	//v.Add("_sort", "ID,desc")
	//v.Add("_sort", "UUID,desc")
	//v.Add("_page", "2")
	//v.Add("_pageSize", "10")
	//v.Add("_select", "ID,UUID")
	//v.Add("q", "Brand.ProductModel=PowerEdge R730xd (SKU=NotProvided;ModelName=PowerEdge R730xd)")
	//v.Encode()

	//r, _ := http.NewRequest("GET", "http://localhost:8080?"+v.Encode(), nil)

	type args struct {
		r *http.Request
		t interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantFns []func(db *gorm.DB) *gorm.DB
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				r: r,
				t: &PhysicalMachine{},
			},
			wantFns: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := SetScopesToCtx(context.Background(), r, *s)
				tmpDB := db.Session(&gorm.Session{DryRun: true})
				data := make([]PhysicalMachine, 0)
				tmpDB, err = SetScopes(ctx, tmpDB)
				if err != nil {
					panic(err)
				}
				fmt.Println(tmpDB.Find(&data))
				fmt.Println(tmpDB.Find(&data).Statement.SQL.String())
				//err = tmpDB.Find(&data).Error
				//if err != nil {
				//	panic(err)
				//}
				//b, _ := json.Marshal(data)
				//fmt.Println(string(b))
			},
		)
	}
}
