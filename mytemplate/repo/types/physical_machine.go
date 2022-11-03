package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

var NetworkCardPortNames = []string{
	"nasPort1",
	"nasPort2",
	"businessPort1",
	"businessPort2",
	"managementPort1",
}

type PhysicalMachine struct {
	gorm.Model
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	//品牌信息
	//Brand string `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`
	//型号
	//BrandModel string `gorm:"column:brand_model;notnull;comment:'型号'" json:"brandModel"`

	// 是否采集完成
	IsCollectFinish bool `gorm:"DEFAULT false; comment:'是否采集完成';column:is_collect_finish" json:"isCollectFinish"`

	BrandUUID string `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     Brand  `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

	// 物理机名字
	PhysicalMachineName string `gorm:"column:physical_machine_name;notnull;comment:'物理机名字'" json:"physicalMachineName"`

	// 业务ip
	BusinessIp string `gorm:"column:business_ip;notnull;comment:'业务ip'" json:"businessIp"`
	// 业务子网掩码
	BusinessSubnetMask string `gorm:"column:business_subnet_mask;notnull;comment:'业务子网掩码'" json:"businessSubnetMask"`
	// nasIp
	NasIp string `gorm:"column:nas_ip;notnull;comment:'nasIp'" json:"nasIp"`
	// nas子网掩码
	NasSubnetMask string `gorm:"column:nas_subnet_mask;notnull;comment:'nas子网掩码'" json:"nasSubnetMask"`
	// IPMIIP
	IpmiIp string `gorm:"column:ipmi_ip;notnull;comment:'ipmiIp'" json:"ipmiIp"`
	// ipmi子网掩码
	IpmiSubnetMask string `gorm:"column:ipmi_subnet_mask;notnull;comment:'ipmi子网掩码'" json:"ipmiSubnetMask"`
	// sn
	Sn string `gorm:"column:sn;notnull;comment:'sn'" json:"sn"`

	// cpu品牌uuid
	//CpuBrandUUID string `gorm:"column:cpu_brand_uuid;notnull;comment:'cpu品牌uuid'" json:"cpuBrandUUID"`
	//CpuBrand     Brand  `gorm:"foreignKey:CpuBrandUUID;references:UUID" json:"cpuBrand"`
	//CpuCount     int    `gorm:"column:cpu_count;notnull;comment:'cpu个数'" json:"cpuCount"`

	// cpu品牌
	//CpuBrand string `gorm:"column:cpu_brand;notnull;comment:'cpu品牌'" json:"cpuBrand"`
	// cpu型号
	//CpuModel string `gorm:"column:cpu_model;notnull;comment:'cpu型号'" json:"cpuModel"`
	// cpu个数
	// 内存大小
	//MemoryBrandUUID string `gorm:"column:memory_brand_uuid;notnull;comment:'内存品牌uuid'" json:"memoryBrandUUID"`
	//MemoryBrand     Brand  `gorm:"foreignKey:MemoryBrandUUID;references:UUID" json:"memoryBrand"`
	//MemoryCount     int    `gorm:"column:memory_count;notnull;comment:'内存数量'" json:"memoryCount"`

	//DiskBrandUUID string `gorm:"column:disk_brand_uuid;notnull;comment:'硬盘品牌uuid'" json:"diskBrandUUID"`
	//DiskBrand     Brand  `gorm:"foreignKey:DiskBrandUUID;references:UUID" json:"diskBrand"`
	//DiskCount     int    `gorm:"column:disk_count;notnull;comment:'硬盘数量'" json:"diskCount"`

	// 网卡品牌UUID
	//NetworkCardBrandUUID string `gorm:"column:network_card_brand_uuid;notnull;comment:'网卡品牌uuid'" json:"networkCardBrandUUID"`
	//NetworkCardBrand     Brand  `gorm:"foreignKey:NetworkCardBrandUUID;references:UUID" json:"networkCardBrand"`
	//NetworkCardCount     int    `gorm:"column:network_card_count;notnull;comment:'网卡数量'" json:"networkCardCount"`

	//MemorySize int `gorm:"column:memory_size;notnull;comment:'内存大小'" json:"memorySize"`
	// 内存数量
	// 硬盘大小
	//DiskSize int `gorm:"column:disk_size;notnull;comment:'硬盘大小'" json:"diskSize"`
	// 硬盘数量
	//DiskCount int `gorm:"column:disk_count;notnull;comment:'硬盘数量'" json:"diskCount"`
	// 网卡数量
	//NetworkCardCount int `gorm:"column:network_card_count;notnull;comment:'网卡数量'" json:"networkCardCount"`
	// 操作系统品牌UUID
	OsBrandUUID string `gorm:"column:os_brand_uuid;notnull;comment:'操作系统品牌uuid'" json:"osBrandUUID"`
	OsBrand     Brand  `gorm:"foreignKey:OsBrandUUID;references:UUID" json:"osBrand"`

	// 供应商UUID
	SupplierUUID string   `gorm:"column:supplier_uuid;notnull;comment:'供应商uuid'" json:"supplierUUID"`
	Supplier     Supplier `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

	// 跳板机登录端口
	JumpPort int `gorm:"column:jump_port;notnull;default:0;comment:'跳板机登录端口'" json:"jumpPort"`
	// 环境  1:生产 2:测试 3:开发
	Env int `gorm:"column:env;notnull;default:1;comment:'环境'" json:"env"`
	// 状态  1:使用中 2：故障 3：关机 4: 进入库房
	Status int `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`

	// 维护状态 1:正常 2:维护中
	MaintainStatus int `gorm:"column:maintain_status;notnull;default:1;comment:'维护状态'" json:"maintainStatus"`

	// 备注
	Remark string `gorm:"column:remark;null;comment:'备注'" json:"remark"`

	NamespaceID int `gorm:"column:namespace_id;notnull;comment:'项目ID'" json:"namespaceId"`

	NamespaceService TblServicetree `gorm:"->;foreignKey:NamespaceID;references:Pri" json:"namespaceService"`

	NameID int `gorm:"column:name_id;notnull;comment:'服务名称'" json:"nameId"`

	NameService TblServicetree `gorm:"->;foreignKey:NameID;references:Pri" json:"nameService"`

	// 维保起始时间
	MaintenanceStartTime string `gorm:"column:maintenance_start_time;notnull;comment:'维保起始时间'" json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `gorm:"column:expire_time;notnull;comment:'维保到期时间'" json:"maintenanceExpireTime"`

	// 负责人
	Owners []string `gorm:"column:owners;serializer:json;notnull;comment:'负责人'" json:"owners"`

	// 关联的设备
	Device Device `gorm:"polymorphic:Device;foreignKey:UUID" json:"device"`

	// 备件关联
	SparePart SparePart `gorm:"polymorphic:Owner;foreignKey:UUID" json:"sparePart"`

	// 配件关联
	Accessories []Accessories `gorm:"polymorphic:Owner;foreignKey:UUID" json:"accessories"`
	// 关联库房
	//Warehouse Warehouse `gorm:"polymorphic:Warehouse" json:"warehouse"`

	// 维护信息
	Maintenances []Maintenance `gorm:"foreignKey:SourceID;references:UUID" json:"maintain"`
}

func (physicalMachine *PhysicalMachine) BeforeCreate(tx *gorm.DB) (err error) {
	if physicalMachine.UUID == "" {
		physicalMachine.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

//func (PhysicalMachine *PhysicalMachine) BeforeUpdate(tx *gorm.DB) (err error) {
//	if PhysicalMachine.Status == 4 {
//		tx.
//	}
//}

func (PhysicalMachine) TableName() string {
	return "assets_physical_machine"
}
