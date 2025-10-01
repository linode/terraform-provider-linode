package helper

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

func ParseRangeOptionalAddress(cidr string) (*netip.Addr, int, error) {
	cidrSplit := strings.Split(cidr, "/")
	if len(cidrSplit) != 2 {
		return nil, 0, fmt.Errorf("malformed CIDR: %s", cidr)
	}

	ipStr, prefixStr := cidrSplit[0], cidrSplit[1]

	if prefixStr == "" {
		return nil, 0, fmt.Errorf("expected non-empty prefix: %s", cidr)
	}

	prefix, err := strconv.Atoi(prefixStr)
	if err != nil {
		return nil, 0, fmt.Errorf("malformed prefix: %w", err)
	}

	// No address is specified
	if ipStr == "" {
		return nil, prefix, nil
	}

	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse CIDR: %w", err)
	}

	return &ip, prefix, nil
}
