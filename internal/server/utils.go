package server

import (
	"strconv"
)

func stringToInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	return i, err
}
