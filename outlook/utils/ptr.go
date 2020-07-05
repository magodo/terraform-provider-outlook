package utils

import "reflect"

func Bool(input bool) *bool {
	return &input
}

func Int(input int) *int {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func Float(input float64) *float64 {
	return &input
}

func String(input string) *string {
	return &input
}

// ToPtr create a new object from the passed in "obj" and return its address back.
func ToPtr(obj interface{}) interface{} {
	v := reflect.ValueOf(obj)
	vp := reflect.New(v.Type())
	vp.Elem().Set(v)
	return vp.Interface()
}

// ToPtrOrNil is similar to ToPtr, except it returns nil if "value" is of "zero" value
func ToPtrOrNil(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface()) {
		return reflect.Zero(reflect.New(v.Type()).Type()).Interface()
	}
	return ToPtr(value)
}

// SafeDeref returns the value pointed to by input pointer. If the pointer is null,
// it returns the "zero" value of the value the input pointer pointed to.
func SafeDeref(ptr interface{}) interface{} {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic("Invalid input: input is not a pointer")
	}
	if v.IsNil() {
		uzero := underlyingZeroValue(v.Type())
		// construct the right type
		rv := reflect.New(v.Type().Elem())
		rv.Elem().Set(uzero)
		return rv.Elem().Interface()
	}
	return v.Elem().Interface()
}
