package di

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/mix-go/xdi"
	"github.com/mix-go/xutil/xenv"
	"xorm.io/xorm"
)

func init() {
	obj := xdi.Object{
		Name: "xorm",
		New: func() (i interface{}, e error) {
			return xorm.NewEngine("mysql", xenv.Getenv("DATABASE_DSN").String())
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}

func Xorm() (db *xorm.Engine) {
	if err := xdi.Populate("xorm", &db); err != nil {
		panic(err)
	}
	return
}
