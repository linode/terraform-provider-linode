package helper

import (
	"context"

	"github.com/linode/linodego"
)

// GetLatestEvent returns the latest Linode event with the given arguments.
func GetLatestEvent(ctx context.Context, client *linodego.Client,
	entityID int, entityType linodego.EntityType, action linodego.EventAction,
) (*linodego.Event, error) {
	filter := linodego.Filter{
		Order:   linodego.Descending,
		OrderBy: "created",
	}
	filter.AddField(linodego.Eq, "action", action)
	filter.AddField(linodego.Eq, "entity.id", entityID)
	filter.AddField(linodego.Eq, "entity.type", entityType)

	filterStr, err := filter.MarshalJSON()
	if err != nil {
		return nil, err
	}

	listOptions := linodego.ListOptions{
		PageOptions: &linodego.PageOptions{Page: 1},
		PageSize:    25,
		Filter:      string(filterStr),
	}

	events, err := client.ListEvents(ctx, &listOptions)
	if err != nil {
		return nil, err
	}

	if len(events) < 1 {
		return nil, nil
	}

	return &events[0], nil
}
