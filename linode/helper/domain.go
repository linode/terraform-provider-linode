package helper

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DomainSecondsDiffSuppressor() schema.SchemaDiffSuppressFunc {
	accepted := []int{
		300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, 2419200,
	}

	rounder := func(n int) int {
		if n == 0 {
			return 0
		}

		for _, value := range accepted {
			if n <= value {
				return value
			}
		}
		return accepted[len(accepted)-1]
	}

	return func(k, provisioned, declared string, d *schema.ResourceData) bool {
		provisionedSec, _ := strconv.Atoi(provisioned)
		declaredSec, _ := strconv.Atoi(declared)
		return rounder(declaredSec) == provisionedSec
	}
}
