package util

import (
	"reflect"

	"github.com/jhseong7/ecl"
)

var (
	log = ecl.NewLogger(ecl.LoggerOption{
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

// Version of DeriveTypeFromInstantiator that returns multiple return types
// This returns all the return types of the instantiator function.
func DeriveTypeListFromInstantiator(instantiator interface{}) ([]reflect.Type, bool) {
	funcType := reflect.TypeOf(instantiator)

	if funcType.Kind() != reflect.Func {
		log.Panicf("Instantiator is not a function: %s", funcType.String())
		return nil, false
	}

	// Check if the return type exists
	if funcType.NumOut() == 0 {
		log.Panicf("Instantiator must have at least one return value: %s", funcType.String())
		return nil, false
	}

	// Get the return types
	returnTypes := make([]reflect.Type, funcType.NumOut())
	for i := 0; i < funcType.NumOut(); i++ {
		returnTypes[i] = funcType.Out(i)
	}

	return returnTypes, true
}

// Retrive the input types of the instantiator function.
func DeriveInputTypesFromInstantiator(instantiator interface{}) ([]reflect.Type, bool) {
	funcType := reflect.TypeOf(instantiator)

	if funcType.Kind() != reflect.Func {
		log.Panicf("Instantiator is not a function: %s", funcType.String())
		return nil, false
	}

	inputTypes := make([]reflect.Type, funcType.NumIn())

	for i := 0; i < funcType.NumIn(); i++ {
		inputTypes[i] = funcType.In(i)
	}

	return inputTypes, true
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
