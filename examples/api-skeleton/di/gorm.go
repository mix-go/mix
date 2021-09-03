package di

import (
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xdi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	obj := xdi.Object{
		Name: "gorm",
		New: func() (i interface{}, e error) {
			return gorm.Open(mysql.Open(dotenv.Getenv("DATABASE_DSN").String()))
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}

func Gorm() (db *gorm.DB) {
	if err := xdi.Populate("gorm", &db); err != nil {
		panic(err)
	}
	return
}
