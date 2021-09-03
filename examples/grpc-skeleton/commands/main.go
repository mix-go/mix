package commands

import (
	"github.com/mix-go/xcli"
)

var Commands = []*xcli.Command{
	{
		Name:  "grpc:server",
		Short: "gRPC server demo",
		Options: []*xcli.Option{
			{
				Names: []string{"d", "daemon"},
				Usage: "Run in the background",
			},
		},
		RunI: &GrpcServerCommand{},
	},
	{
		Name:  "grpc:client",
		Short: "gRPC client demo",
		RunI:  &GrpcClientCommand{},
	},
}
