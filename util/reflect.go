package util

import (
	"reflect"
)

func MergeStruct(base any, option any) error {
	if option == nil || reflect.ValueOf(option).IsNil() {
		return nil
	}
	optionVal := reflect.ValueOf(option).Elem()
	defaultVal := reflect.ValueOf(base).Elem()

	for i := 0; i < optionVal.NumField(); i++ {
		if !optionVal.Field(i).IsZero() {
			defaultVal.Field(i).Set(optionVal.Field(i))
		}
	}

	return nil
}
