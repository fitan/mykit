package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type Accessories struct {
	gorm.Model

	// 配置UUID
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	// 采购批次
	Batch string `gorm:"column:batch;notnull;comment:'采购批次'" json:"batch"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	// sn号
	Sn string `gorm:"column:sn;notnull;comment:'sn号'" json:"sn"`

	OwnerID   string
	OwnerType string

	// 物理机
	PhysicalMachine *PhysicalMachine `gorm:"foreignKey:OwnerID;references:UUID" json:"physicalMachine"`

	// 网络设备
	NetworkDevice *NetworkDevice `gorm:"foreignKey:OwnerID;references:UUID" json:"networkDevice"`

	// 备件关联
	SparePart SparePart `gorm:"polymorphic:Owner;foreignKey:UUID" json:"sparePart"`

	WarehouseUUID string    `gorm:"column:warehouse_uuid;notnull;comment:'仓库uuid'" json:"warehouseUUID"`
	Warehouse     Warehouse `gorm:"foreignKey:WarehouseUUID;references:UUID" json:"warehouse"`

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌UUID'" json:"brandUuid"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

	SupplierUUID string   `gorm:"column:supplier_uuid;notnull;comment:'供应商UUID'" json:"supplierUuid"`
	Supplier     Supplier `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

	// 数据来源
	// 1: 自动采集 2: 手动录入 3: 导入
	DataSource int `gorm:"column:data_source;notnull;comment:'数据来源'" json:"dataSource"`
}

func (Accessories) TableName() string {
	return "assets_accessories"
}

func (a *Accessories) BeforeCreate(tx *gorm.DB) (err error) {
	if a.UUID == "" {
		a.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}
