package utils

import "strconv"

func ParseInt(s string) (int64, error) {
	i, err := strconv.Atoi(s)
	return int64(i), err
}
