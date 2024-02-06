package nodepool

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseNodePoolID(id string) (clusterID int, poolID int, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		clusterID, err = strconv.Atoi(parts[0])
		if err != nil {
			return
		}
		poolID, err = strconv.Atoi(parts[1])
	} else {
		err = fmt.Errorf("error parsing pool ID %s, requires two parts separated by colon, eg clusterID:nodeID", id)
	}
	return
}
