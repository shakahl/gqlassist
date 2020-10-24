package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func FormatFileSize(value interface{}) string {
	defer recovery()

	var size float64

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		size = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		size = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		size = v.Float()
	default:
		return ""
	}

	var KB float64 = 1 << 10
	var MB float64 = 1 << 20
	var GB float64 = 1 << 30
	var TB float64 = 1 << 40
	var PB float64 = 1 << 50

	sizeFormat := func(filesize float64, suffix string) string {
		return strings.Replace(fmt.Sprintf("%.1f %s", filesize, suffix), ".0", "", -1)
	}

	var result string
	if size < KB {
		result = sizeFormat(size, "bytes")
	} else if size < MB {
		result = sizeFormat(size/KB, "KB")
	} else if size < GB {
		result = sizeFormat(size/MB, "MB")
	} else if size < TB {
		result = sizeFormat(size/GB, "GB")
	} else if size < PB {
		result = sizeFormat(size/TB, "TB")
	} else {
		result = sizeFormat(size/PB, "PB")
	}

	return result
}
