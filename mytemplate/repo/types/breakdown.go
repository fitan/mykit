package types

import "gorm.io/gorm"

type Breakdown struct {
	gorm.Model
	UUID string `gorm:"column:uuid;notnull;comment:'uuid'" json:"uuid"`
	//

}
