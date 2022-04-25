package helper

import (
	"context"

	"github.com/linode/linodego"
)

func ResolveValidDBEngine(
	ctx context.Context, client linodego.Client, engine string) (*linodego.DatabaseEngine, error) {
	filter := linodego.Filter{}
	filter.AddField(linodego.Eq, "engine", engine)

	filterBytes, err := filter.MarshalJSON()
	if err != nil {
		return nil, err
	}

	engines, err := client.ListDatabaseEngines(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return nil, err
	}

	return &engines[0], nil
}
