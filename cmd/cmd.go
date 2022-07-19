package main

import (
	"context"

	operator "github.com/jodydadescott/example1/operator"
)

func main() {

	err := run()
	if err != nil {
		panic(err)
	}

}

func run() error {

	ctx := context.Background()

	config := &operator.Config{}
	config.PrismaAPI = "https://api.east-01.network.prismacloud.io"
	config.PrismaLabel = "my-label"
	config.PrismaNamespace = "my-namespace"

	op, err := operator.NewOperator(ctx, config)
	if err != nil {
		return err
	}

	return op.Run(ctx)
}
