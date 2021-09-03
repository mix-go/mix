package configor

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/mix-go/api-skeleton/config"
	"github.com/mix-go/xcli/argv"
)

func init() {
	// Conf support YAML, JSON, TOML, Shell Environment
	if err := configor.Load(&config.Config, fmt.Sprintf("%s/../conf/config.yml", argv.Program().Dir)); err != nil {
		panic(err)
	}
}
