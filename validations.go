package validations

import (
	"reflect"
	"strconv"
)

func DoValidations[T interface{}]() map[string]func(T, reflect.StructField, int) (bool, bool) {
	return map[string]func(T, reflect.StructField, int) (bool, bool){
		reflect.String.String(): func(object T, field reflect.StructField, i int) (isValidationPresent bool, isValidated bool) {
			isValidated = true
			isValidationPresent = false
			if field.Tag.Get("required") == "true" {
				isValidated = reflect.ValueOf(&object).Elem().Field(i).String() != ""
				isValidationPresent = false
			}
			return
		},
		reflect.Int.String(): func(object T, field reflect.StructField, i int) (isValidationPresent bool, isValidated bool) {
			isValidated = true
			isValidationPresent = false
			if field.Tag.Get("required") == "true" {
				isValidated = reflect.ValueOf(&object).Elem().Field(i).Int() != 0
				isValidationPresent = true
			}
			if min, error := strconv.ParseInt(field.Tag.Get("min"), 10, 32); error == nil {
				isValidated = reflect.ValueOf(&object).Elem().Field(i).Int() >= min
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
		doValidation := DoValidations[T]()[field.Type.Kind().String()]
		if isValidationPresent, isValidated := doValidation(data, field, i); isValidationPresent {
			return isValidated
		}
	}
	return true
}
