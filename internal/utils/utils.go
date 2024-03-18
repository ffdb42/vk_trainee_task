package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseID(uri string) (int, error) {
	uri = strings.TrimRight(uri, "/")
	strs := strings.Split(uri, "/")
	if len(strs) != 3 {
		return -1, nil
	}

	id, err := strconv.Atoi(strs[2])
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}
