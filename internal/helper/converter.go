// internal/helper/converter.go
package helper

import (
	"strconv"
)

// StringToIntPtr converts string to *int, returns nil if string is empty or invalid
func StringToIntPtr(s string) *int {
	if s == "" {
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &val
}

// StringToInt converts string to int, returns 0 if string is empty or invalid
func StringToInt(s string) int {
	if s == "" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// IntPtrToString converts *int to string, returns empty string if nil
func IntPtrToString(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}

// StringToFloat64Ptr converts string to *float64, returns nil if string is empty or invalid
func StringToFloat64Ptr(s string) *float64 {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &val
}

// StringToUintPtr converts string to *uint, returns nil if string is empty or invalid
func StringToUintPtr(s string) *uint {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	u := uint(val)
	return &u
}