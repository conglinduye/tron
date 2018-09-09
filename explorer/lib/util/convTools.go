package util

import (
	"fmt"
	"strconv"
)

//[ uint...
func Str2Uint64(str string) uint64 {
	ret, _ := strconv.ParseUint(str, 10, 64)
	return ret
}

func Str2Uint32(str string) uint32 {
	ret, _ := strconv.ParseUint(str, 10, 32)
	return uint32(ret)
}

func Str2Uint16(str string) uint16 {
	ret, _ := strconv.ParseUint(str, 10, 16)
	return uint16(ret)
}

func Str2Uint8(str string) uint8 {
	ret, _ := strconv.ParseUint(str, 10, 8)
	return uint8(ret)
}

func Str2UInt(str string) uint {
	ret, _ := strconv.ParseUint(str, 10, 64)
	return uint(ret)
}

//]

//[ int...
func Str2Int64(str string) int64 {
	ret, _ := strconv.ParseInt(str, 10, 64)
	return ret
}

func Str2Int32(str string) int32 {
	ret, _ := strconv.ParseInt(str, 10, 32)
	return int32(ret)
}

func Str2Int16(str string) int16 {
	ret, _ := strconv.ParseInt(str, 10, 16)
	return int16(ret)
}

func Str2Int8(str string) int8 {
	ret, _ := strconv.ParseInt(str, 10, 8)
	return int8(ret)
}

func Str2Int(str string) int {
	ret, _ := strconv.ParseInt(str, 10, 64)
	return int(ret)
}

//]

//[ float...
func Str2Float64(str string) float64 {
	ret, _ := strconv.ParseFloat(str, 64)
	return ret
}

// Str2Float32 ...
func Str2Float32(str string) float32 {
	ret, _ := strconv.ParseFloat(str, 32)
	return float32(ret)
}

//]

//ToStr [int to str
func ToStr(v interface{}) string {
	ret := fmt.Sprintf("%v", v)
	return ret
}

//]

// ToInt64 return int64 type value from interface, if input is int/uint type, else return 0
func ToInt64(intValue interface{}) int64 {
	ret := int64(0)
	switch v := intValue.(type) {
	case int:
		ret = int64(v)
	case int8:
		ret = int64(v)
	case int16:
		ret = int64(v)
	case int32:
		ret = int64(v)
	case int64:
		ret = int64(v)
	case uint:
		ret = int64(v)
	case uint8:
		ret = int64(v)
	case uint16:
		ret = int64(v)
	case uint32:
		ret = int64(v)
	case uint64:
		ret = int64(v)
	case float32:
		ret = int64(v)
	case float64:
		ret = int64(v)
	default:
		ret = Str2Int64(fmt.Sprintf("%v", v))
	}
	return ret
}
