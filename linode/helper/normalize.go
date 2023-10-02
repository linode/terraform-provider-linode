package helper

import (
	"net"
)

func CompareIPv6Ranges(i, v string) (bool, error) {
	ipi, ipneti, err := net.ParseCIDR(i)
	if err != nil {
		return false, err
	}

	ipv, ipnetv, err := net.ParseCIDR(v)
	if err != nil {
		return false, err
	}

	return ipi.Equal(ipv) && ipneti.Mask.String() == ipnetv.Mask.String(), nil
}
