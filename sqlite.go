package main

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite" // 替换官方驱动
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type XS struct {
	Title   string `gorm:"column:Title" json:"title"` // 设置 Title 为主键
	Chapter string `gorm:"column:Chapter" json:"chapter"`
	Content string `gorm:"column:Content" json:"content"`
	Prompt  string `gorm:"column:Prompt" json:"prompt"`
	IDindex int    `gorm:"column:IDindex" json:"idindex"` // 添加自增ID作为主键
}

var (
	db  *gorm.DB
	err error
)

func init() {

	// 获取当前可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		panic("failed to get executable path")
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "aixiaoshuo.db")

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&XS{}) //AutoMigrate 方法会根据 InitData 结构体自动创建 plcinit 表，如果表已经存在，则会更新表的结构以匹配结构体。

}

func insertXS(xs XS) error {
	return db.Create(&xs).Error
}
func getXs(title string) ([]XS, error) {
	var xsList []XS

	// 正确写法：先指定条件和排序，最后用 Find 执行
	result := db.Where("Title = ?", title).Order("IDindex ASC").Find(&xsList)

	return xsList, result.Error
}
