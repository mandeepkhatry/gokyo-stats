package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetCSVFileName(file string) string {
	d := strings.Split(file, ".")
	return fmt.Sprintf("%s.csv", d[0])
}

func FromToFilter(data string) string {
	res := strings.Contains(data, "_)")
	if res {
		return strings.Split(data, "_")[0]
	}
	return data
}

func GetStringFromInt(data int) string {
	return strconv.Itoa(data)
}

func GetStringFromFloat(data float64) string {
	return fmt.Sprintf("%f", data)
}
