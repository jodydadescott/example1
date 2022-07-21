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

	config := operator.NewConfig().
		SetPrismaAPI("https://api.east-01.network.prismacloud.io").
		SetPrismaLabel("my-label").
		SetPrismaNamespace("my-namespace").
		AddLabelSelectors("node-role.kubernetes.io/master=").
		AddLabelSelectors("node-role.kubernetes.io/infra=")

	op, err := operator.NewOperator(ctx, config)
	if err != nil {
		return err
	}

	return op.Run(ctx)
}
