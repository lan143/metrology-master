package pulsar

import (
	"strconv"
	"strings"
)

func ParseAddress(address string) ([4]byte, error) {
	if strings.HasPrefix(address, "0x") {
		address = address[2:]
	}

	add, err := strconv.ParseUint(address, 16, 32)
	if err != nil {
		return [4]byte{}, err
	}

	return [4]byte{
		byte(add >> 24),
		byte(add >> 16),
		byte(add >> 8),
		byte(add),
	}, nil
}
