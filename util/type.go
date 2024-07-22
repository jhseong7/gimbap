package util

import (
	"reflect"

	logger "github.com/jhseong7/ecl"
)

var (
	log = logger.NewLogger(logger.LoggerOption{
		Name: "UtilLogger",
	})
)

// Retrieve the type of the return value of the instantiator function.
// Also checks if the instantiator is a function with a single return value.
func DeriveTypeFromInstantiator(instantiator interface{}) (reflect.Type, bool) {
	funcType := reflect.TypeOf(instantiator)

	if funcType.Kind() != reflect.Func {
		log.Panicf("Instantiator is not a function: %s", funcType.String())
		return nil, false
	}

	// Check if the return type exists
	if funcType.NumOut() == 0 || funcType.NumOut() > 1 {
		log.Panicf("Instantiator must have a single return value: %s", funcType.String())
		return nil, false
	}

	return funcType.Out(0), true
}

// Remove the pointer from the type. (recursive)
func UnravelPointerType(t reflect.Type, pointerLevel int) (reflect.Type, int) {
	var elem reflect.Type
	if t.Kind() == reflect.Ptr {
		elem = t.Elem()
		pointerLevel++
	} else {
		elem = t
	}

	// If the type of the element is a pointer, unravel it recursively.
	if elem.Kind() == reflect.Ptr {
		return UnravelPointerType(elem, pointerLevel+1)
	}

	return elem, pointerLevel
}

// Get the full name of the type including the package path.
//
// e.g.) "github.com/jhseong7/ecl.Logger"
//
// If the given type is a nested pointer, the name will be prefixed with "*".
// e.g.) "*github.com/jhseong7/ecl.Logger"
func GetFullNameOfType(t reflect.Type) string {
	elem, ptrLevel := UnravelPointerType(t, 0)

	fullName := elem.PkgPath() + "." + elem.Name()

	for i := 0; i < ptrLevel; i++ {
		fullName = "*" + fullName
	}

	return fullName
}
