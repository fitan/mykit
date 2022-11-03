package pm

import (
	types "github.com/fitan/mykit/mytemplate/repo/types"
	pm "github.com/fitan/mykit/mytemplate/services/pm"
)

func listDTO(src []types.PhysicalMachine) (dest []ListResponse) {
	dest = listDTOObj{}.Copy(src)
	return
}

type listDTOObj struct{}

func (d listDTOObj) Copy(src []types.PhysicalMachine) (dest []ListResponse) {
	// basic =
	// slice =
	dest = make([]pm.ListResponse, 0, len(src))
	for i := 0; i < len(src); i++ {
		dest[i] = d.typesPhysicalMachineToPmListResponse(src[i])
	}
	// map =
	// pointer =
	return
}
func (d listDTOObj) typesPhysicalMachineToPmListResponse(src types.PhysicalMachine) (dest pm.ListResponse) {
	// basic =
	// slice =
	/*
	   负责人
	*/
	dest.Owners = src.Owners
	dest.Device.NetworkToDevices = make([]pm.DevicePort, 0, len(src.Device.NetworkToDevices))
	for i := 0; i < len(src.Device.NetworkToDevices); i++ {
		dest.Device.NetworkToDevices[i] = d.typesDevicePortToPmDevicePort(src.Device.NetworkToDevices[i])
	}
	dest.Device.NetworkFromDevices = make([]pm.DevicePort, 0, len(src.Device.NetworkFromDevices))
	for i := 0; i < len(src.Device.NetworkFromDevices); i++ {
		dest.Device.NetworkFromDevices[i] = d.typesDevicePortToPmDevicePort(src.Device.NetworkFromDevices[i])
	}
	dest.Device.Cabinet.ManagerEmail = src.Device.Cabinet.ManagerEmail
	dest.Device.Cabinet.Devices = make([]pm.Device, 0, len(src.Device.Cabinet.Devices))
	for i := 0; i < len(src.Device.Cabinet.Devices); i++ {
		dest.Device.Cabinet.Devices[i] = d.typesDeviceToPmDevice(src.Device.Cabinet.Devices[i])
	}
	/*
	   维护信息
	*/
	dest.Maintenances = make([]struct {
		SourceID        string "gorm:\"column:source_id;notnull;comment:'来源id'\" json:\"sourceId\""
		MaintenanceType string "gorm:\"column:maintenance_type;notnull;comment:'维护类型'\" json:\"maintenanceType\""
		Title           string "gorm:\"column:title;notnull;comment:'维护标题'\" json:\"title\""
		MaintenanceInfo string "gorm:\"column:maintenance_info;notnull;comment:'维护信息'\" json:\"maintenanceInfo\""
		MaintenanceUser string "gorm:\"column:maintenance_user;notnull;comment:'维护人员'\" json:\"maintenanceUser\""
	}, 0, len(src.Maintenances))
	for i := 0; i < len(src.Maintenances); i++ {
		dest.Maintenances[i] = d.typesMaintenanceToXstruct(src.Maintenances[i])
	}
	// map =
	// pointer =
	dest.Brand.ProductType = src.Brand.ProductType
	/*
	   业务ip
	*/
	dest.BusinessIp = src.BusinessIp
	dest.IsCollectFinish = src.IsCollectFinish
	dest.Brand.ProductModel = src.Brand.ProductModel
	dest.Brand.ProductParam.Cpu = src.Brand.ProductParam.Cpu
	/*
	   维护状态 1:正常 2:维护中
	*/
	dest.MaintainStatus = src.MaintainStatus
	/*
	   维保起始时间
	*/
	dest.MaintenanceStartTime = src.MaintenanceStartTime
	dest.Device.Cabinet.Row = src.Device.Cabinet.Row
	dest.BrandUUID = src.BrandUUID
	dest.Brand.ProductParam.Memory = src.Brand.ProductParam.Memory
	/*
	   nas子网掩码
	*/
	dest.NasSubnetMask = src.NasSubnetMask
	/*
	   ipmi子网掩码
	*/
	dest.IpmiSubnetMask = src.IpmiSubnetMask
	/*
	   sn
	*/
	dest.Sn = src.Sn
	/*
	   供应商UUID
	*/
	dest.SupplierUUID = src.SupplierUUID
	dest.Supplier.Name = src.Supplier.Name
	/*
	   状态  1:使用中 2：故障 3：关机 4: 进入库房
	*/
	dest.Status = src.Status
	/*
	   维保到期时间
	*/
	dest.MaintenanceExpireTime = src.MaintenanceExpireTime
	dest.Device.EndLayer = src.Device.EndLayer
	dest.Device.CabinetID = src.Device.CabinetID
	dest.Device.Cabinet.ComputerRoomID = src.Device.Cabinet.ComputerRoomID
	dest.Brand.Model.ID = src.Model.ID
	dest.Brand.ProductParam.HardDisk = src.Brand.ProductParam.HardDisk
	/*
	   业务子网掩码
	*/
	dest.BusinessSubnetMask = src.BusinessSubnetMask
	dest.Supplier.Address = src.Supplier.Address
	dest.NameID = src.NameID
	dest.Device.DeviceID = src.Device.DeviceID
	dest.Device.StartLayer = src.Device.StartLayer
	dest.Device.Cabinet.Code = src.Device.Cabinet.Code
	/*
	   备注
	*/
	dest.Remark = src.Remark
	/*
	   nasIp
	*/
	dest.NasIp = src.NasIp
	/*
	   IPMIIP
	*/
	dest.IpmiIp = src.IpmiIp
	dest.OsBrandUUID = src.OsBrandUUID
	dest.Device.Cabinet.Column = src.Device.Cabinet.Column
	/*
	   uuid
	*/
	dest.UUID = src.UUID
	dest.Brand.Model.DeletedAt.Valid = src.Model.DeletedAt.Valid
	dest.Brand.Brand = src.Brand.Brand
	/*
	   跳板机登录端口
	*/
	dest.JumpPort = src.JumpPort
	dest.NamespaceID = src.NamespaceID
	dest.Device.DeviceType = src.Device.DeviceType
	dest.Device.Cabinet.Type = src.NamespaceService.Type
	/*
	   物理机名字
	*/
	dest.PhysicalMachineName = src.PhysicalMachineName
	dest.Supplier.Contact = src.Supplier.Contact
	/*
	   环境  1:生产 2:测试 3:开发
	*/
	dest.Env = src.Env
	dest.Device.Cabinet.Size = src.Device.Cabinet.Size
	return
}
func (d listDTOObj) typesDevicePortToPmDevicePort(src types.DevicePort) (dest pm.DevicePort) {
	// basic =
	// slice =
	dest.Device.NetworkToDevices = make([]pm.DevicePort, 0, len(src.Device.NetworkToDevices))
	for i := 0; i < len(src.Device.NetworkToDevices); i++ {
		dest.Device.NetworkToDevices[i] = d.typesDevicePortToPmDevicePort(src.Device.NetworkToDevices[i])
	}
	dest.Device.NetworkFromDevices = make([]pm.DevicePort, 0, len(src.Device.NetworkFromDevices))
	for i := 0; i < len(src.Device.NetworkFromDevices); i++ {
		dest.Device.NetworkFromDevices[i] = d.typesDevicePortToPmDevicePort(src.Device.NetworkFromDevices[i])
	}
	dest.Device.Cabinet.ManagerEmail = src.Device.Cabinet.ManagerEmail
	dest.Device.Cabinet.Devices = make([]pm.Device, 0, len(src.Device.Cabinet.Devices))
	for i := 0; i < len(src.Device.Cabinet.Devices); i++ {
		dest.Device.Cabinet.Devices[i] = d.typesDeviceToPmDevice(src.Device.Cabinet.Devices[i])
	}
	// map =
	// pointer =
	dest.DeviceID = src.DeviceID
	dest.ToDevicePortID = src.ToDevicePortID
	dest.Device.Cabinet.Column = src.Device.Cabinet.Column
	dest.Device.Cabinet.Type = src.Device.Cabinet.Type
	dest.Device.Cabinet.Size = src.Device.Cabinet.Size
	dest.DevicePortName = src.DevicePortName
	if src.ToDevicePort != nil {
		v := d.typesDevicePortToTypesDevicePort(*src.ToDevicePort)
		dest.ToDevicePort = &v
	} else {
		dest.ToDevicePort = src.ToDevicePort
	}
	dest.Device.DeviceType = src.Device.DeviceType
	dest.Device.EndLayer = src.Device.EndLayer
	dest.Device.Cabinet.Name = src.Device.Cabinet.Name
	dest.Device.Cabinet.Code = src.Device.Cabinet.Code
	dest.Device.Cabinet.ComputerRoomID = src.Device.Cabinet.ComputerRoomID
	if src.FromDevicePort != nil {
		v := d.typesDevicePortToTypesDevicePort(*src.FromDevicePort)
		dest.FromDevicePort = &v
	} else {
		dest.FromDevicePort = src.FromDevicePort
	}
	dest.Device.CabinetID = src.Device.CabinetID
	dest.Device.Cabinet.Row = src.Device.Cabinet.Row
	dest.UUID = src.UUID
	dest.ToDeviceID = src.ToDeviceID
	dest.Device.StartLayer = src.Device.StartLayer
	dest.Device.Cabinet.Status = src.Device.Cabinet.Status
	dest.Device.Cabinet.Remark = src.Device.Cabinet.Remark
	return
}
func (d listDTOObj) typesDeviceToPmDevice(src types.Device) (dest pm.Device) {
	// basic =
	// slice =
	dest.Cabinet.ManagerEmail = src.Cabinet.ManagerEmail
	dest.Cabinet.Devices = make([]pm.Device, 0, len(src.Cabinet.Devices))
	for i := 0; i < len(src.Cabinet.Devices); i++ {
		dest.Cabinet.Devices[i] = d.typesDeviceToPmDevice(src.Cabinet.Devices[i])
	}
	dest.NetworkToDevices = make([]pm.DevicePort, 0, len(src.NetworkToDevices))
	for i := 0; i < len(src.NetworkToDevices); i++ {
		dest.NetworkToDevices[i] = d.typesDevicePortToPmDevicePort(src.NetworkToDevices[i])
	}
	dest.NetworkFromDevices = make([]pm.DevicePort, 0, len(src.NetworkFromDevices))
	for i := 0; i < len(src.NetworkFromDevices); i++ {
		dest.NetworkFromDevices[i] = d.typesDevicePortToPmDevicePort(src.NetworkFromDevices[i])
	}
	// map =
	// pointer =
	dest.CabinetID = src.CabinetID
	dest.Cabinet.ComputerRoomID = src.Cabinet.ComputerRoomID
	dest.Cabinet.Row = src.Cabinet.Row
	dest.Cabinet.Column = src.Cabinet.Column
	dest.Cabinet.Remark = src.Cabinet.Remark
	dest.UUID = src.UUID
	dest.DeviceID = src.DeviceID
	dest.DeviceType = src.DeviceType
	dest.Cabinet.Code = src.Cabinet.Code
	dest.Cabinet.Type = src.Cabinet.Type
	dest.Cabinet.Size = src.Cabinet.Size
	dest.Cabinet.Status = src.Cabinet.Status
	dest.StartLayer = src.StartLayer
	dest.EndLayer = src.EndLayer
	dest.Cabinet.Name = src.Cabinet.Name
	return
}
func (d listDTOObj) typesDevicePortToTypesDevicePort(src types.DevicePort) (dest types.DevicePort) {
	// basic =
	// slice =
	dest.Device.NetworkToDevices = make([]pm.DevicePort, 0, len(src.Device.NetworkToDevices))
	for i := 0; i < len(src.Device.NetworkToDevices); i++ {
		dest.Device.NetworkToDevices[i] = d.typesDevicePortToPmDevicePort(src.Device.NetworkToDevices[i])
	}
	dest.Device.NetworkFromDevices = make([]pm.DevicePort, 0, len(src.Device.NetworkFromDevices))
	for i := 0; i < len(src.Device.NetworkFromDevices); i++ {
		dest.Device.NetworkFromDevices[i] = d.typesDevicePortToPmDevicePort(src.Device.NetworkFromDevices[i])
	}
	dest.Device.Cabinet.ManagerEmail = src.Device.Cabinet.ManagerEmail
	dest.Device.Cabinet.Devices = make([]pm.Device, 0, len(src.Device.Cabinet.Devices))
	for i := 0; i < len(src.Device.Cabinet.Devices); i++ {
		dest.Device.Cabinet.Devices[i] = d.typesDeviceToPmDevice(src.Device.Cabinet.Devices[i])
	}
	// map =
	// pointer =
	dest.DevicePortName = src.DevicePortName
	dest.ToDeviceID = src.ToDeviceID
	if src.FromDevicePort != nil {
		v := d.typesDevicePortToTypesDevicePort(*src.FromDevicePort)
		dest.FromDevicePort = &v
	} else {
		dest.FromDevicePort = src.FromDevicePort
	}
	dest.Device.CabinetID = src.Device.CabinetID
	dest.Device.Cabinet.Column = src.Device.Cabinet.Column
	dest.UUID = src.UUID
	dest.ToDevicePortID = src.ToDevicePortID
	dest.Device.StartLayer = src.Device.StartLayer
	dest.Device.Cabinet.Status = src.Device.Cabinet.Status
	dest.DeviceID = src.DeviceID
	dest.Device.Cabinet.Name = src.Device.Cabinet.Name
	dest.Device.Cabinet.Type = src.Device.Cabinet.Type
	dest.Device.Cabinet.Size = src.Device.Cabinet.Size
	dest.Device.Cabinet.Row = src.Device.Cabinet.Row
	dest.Device.Cabinet.ComputerRoomID = src.Device.Cabinet.ComputerRoomID
	if src.ToDevicePort != nil {
		v := d.typesDevicePortToTypesDevicePort(*src.ToDevicePort)
		dest.ToDevicePort = &v
	} else {
		dest.ToDevicePort = src.ToDevicePort
	}
	dest.Device.DeviceType = src.Device.DeviceType
	dest.Device.EndLayer = src.Device.EndLayer
	dest.Device.Cabinet.Code = src.Device.Cabinet.Code
	dest.Device.Cabinet.Remark = src.Device.Cabinet.Remark
	return
}
func (d listDTOObj) typesMaintenanceToXstruct(src types.Maintenance) (dest struct {
	SourceID        string "gorm:\"column:source_id;notnull;comment:'来源id'\" json:\"sourceId\""
	MaintenanceType string "gorm:\"column:maintenance_type;notnull;comment:'维护类型'\" json:\"maintenanceType\""
	Title           string "gorm:\"column:title;notnull;comment:'维护标题'\" json:\"title\""
	MaintenanceInfo string "gorm:\"column:maintenance_info;notnull;comment:'维护信息'\" json:\"maintenanceInfo\""
	MaintenanceUser string "gorm:\"column:maintenance_user;notnull;comment:'维护人员'\" json:\"maintenanceUser\""
}) {
	// basic =
	// slice =
	// map =
	// pointer =
	dest.MaintenanceUser = src.MaintenanceUser
	dest.SourceID = src.SourceID
	dest.MaintenanceType = src.MaintenanceType
	dest.Title = src.Title
	dest.MaintenanceInfo = src.MaintenanceInfo
	return
}
