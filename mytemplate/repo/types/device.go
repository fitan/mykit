package types

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Device struct {
	gorm.Model
	// 设备UUID
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	// 机柜ID
	CabinetID uint `gorm:"column:cabinet_id;notnull;comment:'机柜id'" json:"cabinetId"`
	// 设备id
	DeviceID string `gorm:"column:device_id;notnull;comment:'设备id'" json:"deviceId"`
	// 设备类型
	DeviceType string `gorm:"column:device_type;notnull;comment:'设备类型'" json:"deviceType"`
	// 起始层
	StartLayer int `gorm:"column:start_layer;notnull;comment:'起始层'" json:"startLayer"`
	// 终止层
	EndLayer int `gorm:"column:end_layer;notnull;comment:'终止层'" json:"endLayer"`
	// 网线连接设备
	NetworkToDevices []DevicePort `gorm:"foreignKey:DeviceID;references:DeviceID" json:"networkToDevices"`

	// 别的设备连接至此设备
	NetworkFromDevices []DevicePort `gorm:"foreignKey:ToDeviceID;references:DeviceID" json:"networkFromDevices"`
	// 关联的机柜
	Cabinet Cabinet `gorm:"foreignKey:CabinetID" json:"cabinet"`
}

func (device *Device) BeforeCreate(tx *gorm.DB) (err error) {
	//err = device.CanPutIntoCabinet(tx)
	//if err != nil {
	//	err = errors.Wrap(err, "BeforeCreate.device.CanPutIntoCabinet")
	//	return
	//}
	if device.UUID == "" {
		device.UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return
}

//func (device *Device) BeforeUpdate(tx *gorm.DB) (err error) {
//err = device.CanPutIntoCabinet(tx)
//if err != nil {
//	err = errors.Wrap(err, "BeforeUpdate.device.CanPutIntoCabinet")
//	return
//}
//return
//}

func (device *Device) CanPutIntoCabinet(tx *gorm.DB) (err error) {
	if device.StartLayer > device.EndLayer {
		err = errors.New("起始层不能大于结束层")
		return
	}

	var cabinet Cabinet
	var cabinetID uint
	if device.Cabinet.ID > 0 {
		cabinetID = device.Cabinet.ID
	} else {
		cabinetID = device.CabinetID
	}
	err = tx.Model(&Cabinet{}).Where("id = ?", cabinetID).Preload("ComputerRoom").Preload("Devices", "uuid <> ?", device.UUID).First(&cabinet).Error
	if err != nil {
		err = errors.Wrap(err, "cabinet by uuid")
		return
	}

	if device.EndLayer > cabinet.Size {
		err = errors.New("超出机柜大小 " + strconv.Itoa(cabinet.Size))
		return
	}

	for _, deviceV := range cabinet.Devices {
		if (deviceV.StartLayer == device.StartLayer && deviceV.EndLayer == device.EndLayer) || (deviceV.StartLayer >= device.StartLayer && deviceV.StartLayer <= device.EndLayer) || (deviceV.EndLayer >= device.StartLayer && deviceV.EndLayer <= device.EndLayer) {
			var id string
			err = tx.Table(deviceV.DeviceType).Where("id = ?", deviceV.DeviceID).Pluck("uuid", &id).Error
			if err != nil {
				return
			}
			err = fmt.Errorf("需要更新的设备uuid:%s 起始层或结束层被设备占用，设备类型： %s 设备类型uuid: %s 设备uuid：%s ", device.UUID, deviceV.DeviceType, deviceV.UUID, id)
			return
		}
	}

	return
}

func (Device) TableName() string {
	return "assets_device"
}
