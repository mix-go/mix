package dotenv

import (
	"fmt"
	"github.com/mix-go/xcli/argv"
	"github.com/mix-go/dotenv"
)

func init()  {
	// Env
	if err := dotenv.Load(fmt.Sprintf("%s/../.env", argv.Program().Dir)); err != nil {
		panic(err)
	}
}