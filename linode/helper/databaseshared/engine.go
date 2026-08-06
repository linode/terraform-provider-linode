package databaseshared

import (
	"context"
	"fmt"
	"strings"

	"github.com/linode/linodego/v2"
)

func ResolveValidDBEngine(
	ctx context.Context, client linodego.Client, engine string,
) (*linodego.DatabaseEngine, error) {
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

	if len(engines) < 1 {
		return nil, fmt.Errorf("no db engines were found")
	}

	return &engines[0], nil
}

func CreateLegacyDatabaseEngineSlug(engine, version string) string {
	return fmt.Sprintf("%s/%s", engine, version)
}

func CreateDatabaseEngineSlug(engine, version string) string {
	parts := strings.Split(version, ".")

	switch engine {
	case "mysql":
		if len(parts) >= 2 {
			// When minor part of the engine version is 0, then API accepts major version number only to
			// be used in config (e.g. mysql/8, not mysql/8.0.45). It is not consistent with the API response,
			// which always returns full engine version (mysql/8.0.45). It leads to a discrepancies between
			// config and state refreshed with the API so it needs to be handled
			if parts[1] == "0" {
				return fmt.Sprintf("%s/%s", engine, parts[0])
			}

			return fmt.Sprintf("%s/%s.%s", engine, parts[0], parts[1])
		}

		return fmt.Sprintf("%s/%s", engine, version)

	case "postgresql":
		if len(parts) >= 1 {
			return fmt.Sprintf("%s/%s", engine, parts[0])
		}

		return fmt.Sprintf("%s/%s", engine, version)

	default:
		if len(parts) >= 1 {
			return fmt.Sprintf("%s/%s", engine, parts[0])
		}

		return fmt.Sprintf("%s/%s", engine, version)
	}
}

func ParseDatabaseEngineSlug(engineID string) (string, string, error) {
	components := strings.Split(engineID, "/")
	if len(components) != 2 {
		return "", "", fmt.Errorf("invalid number of components: %d != 2", len(components))
	}

	return components[0], components[1], nil
}
