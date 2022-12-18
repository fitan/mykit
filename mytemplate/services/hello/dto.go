package hello

import gorm "gorm.io/gorm"

func queryDTO(v Query) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("age = ?", v.Age)
		db = db.Where("email = ?", v.Email)
		db = db.Where("id = ?", v.IDIn)
		db = db.Where("time = ?", v.BetweenTime)

		return db
	}
}
