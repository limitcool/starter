package util

import (
	"strconv"
)

// ParseUint 将字符串转换为uint
// 如果转换失败，则返回默认值
func ParseUint(str string, defaultVal uint) uint {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return defaultVal
	}
	return uint(val)
}

// ParseInt 将字符串转换为int
// 如果转换失败，则返回默认值
func ParseInt(str string, defaultVal int) int {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultVal
	}
	return int(val)
}

// ParseFloat 将字符串转换为float64
// 如果转换失败，则返回默认值
func ParseFloat(str string, defaultVal float64) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseBool 将字符串转换为bool
// 如果转换失败，则返回默认值
func ParseBool(str string, defaultVal bool) bool {
	val, err := strconv.ParseBool(str)
	if err != nil {
		return defaultVal
	}
	return val
}
