package main

import (
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	zl, _ := zap.NewDevelopment()
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: glogger.New(goormLoggerFunc(logger.NewZapLogger(zl).Debug), glogger.Config{
			// 慢查询
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// db = db.Debug()

	// 迁移 schema
	// 初始化你的表结构
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整型主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	// WHERE 条件是什么？
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)

	err = db.Raw("SELECT * FROM `products`").First(&product).Error

	// 增删改
	//db.Exec("UPDATE")

	//res := db.Raw("SELECT").Rows()
}

type goormLoggerFunc func(msg string, fields ...logger.Field)

func (g goormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{Key: "args", Val: i})
}
