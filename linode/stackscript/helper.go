package stackscript

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

func UpgradeTimeFormatToRFC3339(oldTime string) (timetypes.RFC3339, error) {
	if newTime, err := time.Parse(time.RFC3339, oldTime); err == nil {
		return timetypes.NewRFC3339TimeValue(newTime), nil
	}

	newTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", oldTime)
	if err != nil {
		return timetypes.NewRFC3339Null(), err
	}

	return timetypes.NewRFC3339TimeValue(newTime), nil
}
