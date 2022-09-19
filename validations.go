package validations

import (
	"reflect"
	"strconv"
)

func doValidations[T interface{}]() map[string]func(T, reflect.StructField, int) (bool, bool) {
	return map[string]func(T, reflect.StructField, int) (bool, bool){
		reflect.String.String(): func(object T, field reflect.StructField, i int) (isValidationPresent bool, isValidated bool) {
			isValidated = true
			isValidationPresent = false
			if field.Tag.Get("required") == "true" {
				isValidated = reflect.ValueOf(&object).Elem().Field(i).String() != ""
				isValidationPresent = true
			}
			return
		},
		reflect.Int.String(): func(object T, field reflect.StructField, i int) (isValidationPresent bool, isValidated bool) {
			isValidated = true
			isValidationPresent = false
			if min, error := strconv.ParseInt(field.Tag.Get("min"), 10, 32); error == nil {
				isValidated = reflect.ValueOf(&object).Elem().Field(i).Int() >= min
				isValidationPresent = true
			}
			if max, error := strconv.ParseInt(field.Tag.Get("max"), 10, 32); error == nil {
				isValidated = reflect.ValueOf(&object).Elem().Field(i).Int() <= max
				isValidationPresent = true
			}
			return
		},
	}
}

func IsValid[T interface{}](data T) bool {
	t := reflect.TypeOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		doValidation := doValidations[T]()[field.Type.Kind().String()]
		/**
		if isValidated is false then return isValidated, otherwise continue validation
		*/
		if isValidationPresent, isValidated := doValidation(data, field, i); isValidationPresent && !isValidated {
			return isValidated
		}
	}
	return true
}
