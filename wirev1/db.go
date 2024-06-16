package wirev1

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("xxx"))
	if err != nil {
		panic(err)
	}
	return db
}
