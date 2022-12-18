package mygorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

type PhysicalMachine struct {
	gorm.Model
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

type Brand struct {
	gorm.Model
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
}

func TestQ(t *testing.T) {
	dsn := "root:123456@tcp(172.29.107.199:3306)/gteml?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DryRun: true})
	if err != nil {
		panic(err)
	}
	fmt.Println(db)

	r, _ := http.NewRequest("GET", "http://localhost:8080?q=Brand.ProductType=1&q=UUID=423424fjdfsdaf234234", nil)

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
				gotFns, err := Q(tt.args.r, tt.args.t)
				if err != nil {
					panic(err)
				}
				tmpDB := db
				fmt.Println(gotFns)
				for _, v := range gotFns {
					tmpDB = v(tmpDB)
				}
				fmt.Println(tmpDB.Find(&PhysicalMachine{}).Statement.SQL.String())
			},
		)
	}
}
