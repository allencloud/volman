package utils

import (
	"strconv"
)

func Convert(in int64) string {
	input := int(in)
	if input < 1024 {
		return strconv.Itoa(input) + "B"
	} else if input < 1024*1024 {
		return strconv.Itoa(input/1024) + "KB"
	} else {
		return strconv.Itoa(input/1024/1024) + "MB"
	}
}
