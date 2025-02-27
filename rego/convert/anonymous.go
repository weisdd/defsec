package convert

import (
	"reflect"
)

var converterInterface = reflect.TypeOf((*Converter)(nil)).Elem()

func anonymousToRego(inputValue reflect.Value) interface{} {

	if inputValue.IsZero() {
		return nil
	}

	if inputValue.Type().Implements(converterInterface) {
		returns := inputValue.MethodByName("ToRego").Call(nil)
		return returns[0].Interface()
	}

	for inputValue.Type().Kind() == reflect.Ptr {
		if inputValue.IsNil() {
			return nil
		}
		inputValue = inputValue.Elem()
	}

	switch kind := inputValue.Type().Kind(); kind {
	case reflect.Struct:
		return StructToRego(inputValue)
	case reflect.Slice:
		return SliceToRego(inputValue)
	}

	return nil
}
