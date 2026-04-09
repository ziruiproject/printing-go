package helper

import (
	"reflect"
	"strconv"
	"strings"
)

func Differ(oldObj, newObj interface{}) interface{} {
	oldVal := reflect.ValueOf(oldObj)
	newVal := reflect.ValueOf(newObj)

	if oldVal.Kind() == reflect.Ptr {
		oldVal = oldVal.Elem()
	}
	if newVal.Kind() == reflect.Ptr {
		newVal = newVal.Elem()
	}

	result := reflect.New(oldVal.Type()).Elem()
	result.Set(oldVal)

	for i := 0; i < oldVal.NumField(); i++ {
		oldFieldVal := oldVal.Field(i)
		newFieldVal := newVal.Field(i)

		if !oldFieldVal.CanInterface() || !newFieldVal.CanInterface() {
			continue
		}

		newValInterface := newFieldVal.Interface()
		oldValInterface := oldFieldVal.Interface()

		if isZeroValue(newValInterface) {
			continue
		}

		if !reflect.DeepEqual(oldValInterface, newValInterface) {
			result.Field(i).Set(newFieldVal)
		}
	}

	return result.Interface()
}

func isZeroValue(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

func ParseUintSlice(s string) []uint {
	var result []uint
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if n, err := strconv.ParseUint(part, 10, 32); err == nil {
			result = append(result, uint(n))
		}
	}
	return result
}
