package dotenv

import (
	"fmt"
	"github.com/mix-go/xcli/argv"
	"github.com/mix-go/xutil/xenv"
)

func init() {
	// Env
	if err := xenv.Load(fmt.Sprintf("%s/../.env", argv.Program().Dir)); err != nil {
		panic(err)
	}
}
