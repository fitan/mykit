package pm

import (
	"github.com/fitan/mykit/mytemplate/repo/types"
)

type ListRequest struct {
	Page  int `json:"page" param:"path,page"`
	Limit int `json:"limit" param:"path,limit"`
}

type DevicePort struct {
	UUID           string      `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	DeviceID       string      `gorm:"column:device_id;notnull;comment:'设备id'" json:"deviceId"`
	DevicePortName string      `gorm:"column:device_port_name;notnull;comment:'设备端口名称'" json:"devicePortName"`
	ToDeviceID     string      `gorm:"column:to_device_id;notnull;comment:'连接至设备id'" json:"toDeviceId"`
	ToDevicePortID string      `gorm:"column:to_device_port_id;notnull;comment:'对端设备端口id'" json:"toDevicePortId"`
	ToDevicePort   *DevicePort `gorm:"foreignKey:ToDevicePortID;references:UUID" json:"toDevicePort"`
	FromDevicePort *DevicePort `gorm:"foreignKey:UUID;references:ToDevicePortID" json:"fromDevicePort"`
	Device         Device      `gorm:"foreignKey:DeviceID;references:DeviceID" json:"sourceDevice"`
}

type Cabinet struct {
	UUID           string   `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	Name           string   `gorm:"column:name;notnull;comment:'名称'" json:"name"`
	Code           string   `gorm:"column:code;notnull;comment:'编号'" json:"code"`
	Type           int      `gorm:"column:type;notnull;default:0;comment:'类型'" json:"type"`
	Size           int      `gorm:"column:size;notnull;default:0;comment:'大小'" json:"size"`
	Row            int      `gorm:"column:row;notnull;default:0;comment:'第几行'" json:"row"`
	Column         int      `gorm:"column:column;notnull;default:0;comment:'第几列'" json:"column"`
	Status         int      `gorm:"column:status;notnull;default:1;comment:'状态'" json:"status"`
	ManagerEmail   []string `gorm:"column:manager_email;serializer:json;notnull;comment:'管理人邮箱'" json:"managerEmail"`
	Remark         string   `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	ComputerRoomID uint     `gorm:"column:computer_room_id;notnull;comment:'机房ID'" json:"computerRoomId"`
	Devices        []Device `gorm:"foreignKey:CabinetID" json:"device"`
}

type Device struct {
	UUID               string       `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	CabinetID          uint         `gorm:"column:cabinet_id;notnull;comment:'机柜id'" json:"cabinetId"`
	DeviceID           string       `gorm:"column:device_id;notnull;comment:'设备id'" json:"deviceId"`
	DeviceType         string       `gorm:"column:device_type;notnull;comment:'设备类型'" json:"deviceType"`
	StartLayer         int          `gorm:"column:start_layer;notnull;comment:'起始层'" json:"startLayer"`
	EndLayer           int          `gorm:"column:end_layer;notnull;comment:'终止层'" json:"endLayer"`
	NetworkToDevices   []DevicePort `gorm:"foreignKey:DeviceID;references:DeviceID" json:"networkToDevices"`
	NetworkFromDevices []DevicePort `gorm:"foreignKey:ToDeviceID;references:DeviceID" json:"networkFromDevices"`
	Cabinet            Cabinet      `gorm:"foreignKey:CabinetID" json:"cabinet"`
}

type ListResponse struct {
	// uuid
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`

	IsCollectFinish bool `gorm:"DEFAULT false; comment:'是否采集完成';column:is_collect_finish" json:"isCollectFinish"`

	BrandUUID string      `gorm:"column:brand_uuid;notnull;comment:'品牌uuid'" json:"brandUUID"`
	Brand     types.Brand `gorm:"foreignKey:BrandUUID;references:UUID" json:"brand"`

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

	OsBrandUUID string `gorm:"column:os_brand_uuid;notnull;comment:'操作系统品牌uuid'" json:"osBrandUUID"`
	OsBrand     Brand  `gorm:"foreignKey:OsBrandUUID;references:UUID" json:"osBrand"`

	// 供应商UUID
	SupplierUUID string `gorm:"column:supplier_uuid;notnull;comment:'供应商uuid'" json:"supplierUUID"`
	Supplier     struct {
		UUID    string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
		Name    string `gorm:"column:name;notnull;comment:'名字'" json:"name"`
		Contact string `gorm:"column:contact;notnull;comment:'联系方式'" json:"contact"`
		Address string `gorm:"column:address;notnull;comment:'地址'" json:"address"`
		Remark  string `gorm:"column:remark;null;comment:'备注'" json:"remark"`
	} `gorm:"foreignKey:SupplierUUID;references:UUID" json:"supplier"`

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

	NameID int `gorm:"column:name_id;notnull;comment:'服务名称'" json:"nameId"`

	// 维保起始时间
	MaintenanceStartTime string `gorm:"column:maintenance_start_time;notnull;comment:'维保起始时间'" json:"maintenanceStartTime"`

	// 维保到期时间
	MaintenanceExpireTime string `gorm:"column:expire_time;notnull;comment:'维保到期时间'" json:"maintenanceExpireTime"`

	// 负责人
	Owners []string `gorm:"column:owners;serializer:json;notnull;comment:'负责人'" json:"owners"`

	// 关联的设备
	Device Device `gorm:"polymorphic:Device;foreignKey:UUID" json:"device"`

	// 关联库房
	//Warehouse Warehouse `gorm:"polymorphic:Warehouse" json:"warehouse"`

	// 维护信息
	Maintenances []struct {
		SourceID        string `gorm:"column:source_id;notnull;comment:'来源id'" json:"sourceId"`
		MaintenanceType string `gorm:"column:maintenance_type;notnull;comment:'维护类型'" json:"maintenanceType"`
		Title           string `gorm:"column:title;notnull;comment:'维护标题'" json:"title"`
		MaintenanceInfo string `gorm:"column:maintenance_info;notnull;comment:'维护信息'" json:"maintenanceInfo"`
		MaintenanceUser string `gorm:"column:maintenance_user;notnull;comment:'维护人员'" json:"maintenanceUser"`
	} `gorm:"foreignKey:SourceID;references:UUID" json:"maintain"`
}

type Brand struct {
	UUID         string             `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	Brand        string             `gorm:"column:brand;notnull;comment:'品牌'" json:"brand"`
	ProductType  string             `gorm:"column:product_type;notnull;comment:'产品类型'" json:"productType"`
	ProductModel string             `gorm:"column:product_model;notnull;comment:'产品型号'" json:"productModel"`
	ProductParam types.ProductParam `gorm:"column:product_param;serializer:json;notnull;default:'{}';comment:'产品参数'" json:"productParam"`
	Remark       string             `gorm:"column:remark;null;comment:'备注'" json:"remark"`
}
