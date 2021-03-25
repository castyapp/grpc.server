package helpers

import (
	"net/url"
	"reflect"
)

func GetValidationErrorsFromGoValidator(errs url.Values) map[string] interface{} {
	errors := map[string] interface{} {}
	for index, err := range errs {
		reflect.ValueOf(errors).SetMapIndex(reflect.ValueOf(index), reflect.ValueOf(err))
	}
	return errors
}
