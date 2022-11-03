package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

// 网络设备
type NetworkDevice struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	// 网络设备名字
	NetworkDeviceName string `gorm:"column:network_device_name;notnull;comment:'网络设备名字'" json:"networkDeviceName"`

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

	// 供应商UUID
	SupplierUUID string   `gorm:"column:supplier_uuid;notnull;comment:'供应商uuid'" json:"supplierUUID"`
	Supplier     Supplier `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

	ManagementIP string `gorm:"column:management_ip;notnull;comment:'管理ip'" json:"managementIp"`
	// sn
	Sn string `gorm:"column:sn;notnull;comment:'sn'" json:"sn"`
	// 端口数量
	PortCount int `gorm:"column:port_count;notnull;default:0;comment:'端口数量'" json:"portCount"`
	// 环境  1:生产 2:测试 3:开发
	Env int `gorm:"column:env;notnull;default:1;comment:'环境'" json:"env"`
	// 状态  1:使用中 2：故障 3：关机 4: 进入库房
	Status int `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`

	// 维护状态 1:正常 2:维护中
	MaintainStatus int `gorm:"column:maintain_status;notnull;default:1;comment:'维护状态'" json:"maintainStatus"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	// 关联的设备
	Device Device `gorm:"polymorphic:Device;foreignKey:UUID" json:"device"`

	// 备件关联
	SparePart SparePart `gorm:"polymorphic:Owner;foreignKey:UUID" json:"sparePart"`

	// 配件关联
	Accessories []Accessories `gorm:"polymorphic:Owner;foreignKey:UUID" json:"accessories"`

	// 关联库房
	//Warehouse Warehouse `gorm:"polymorphic:Warehouse" json:"warehouse"`
}

func (networkDevice *NetworkDevice) BeforeCreate(tx *gorm.DB) (err error) {
	if networkDevice.UUID == "" {
		networkDevice.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

func (NetworkDevice) TableName() string {
	return "assets_network_device"
}
