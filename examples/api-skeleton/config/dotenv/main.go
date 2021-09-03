package dotenv

import (
	"fmt"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xcli/argv"
)

func init() {
	// Env
	if err := dotenv.Load(fmt.Sprintf("%s/../.env", argv.Program().Dir)); err != nil {
		panic(err)
	}
}
