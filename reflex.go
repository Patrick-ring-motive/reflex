package reflex

import (
	"reflect"
)

func TypeOf(typ interface{}) reflect.Type {
	switch v := typ.(type) {
	case reflect.Type:
		return v
	case reflect.Value:
		return v.Type()
	default:
		return reflect.TypeOf(v)
	}
}

func ValueOf(val interface{}) reflect.Value {
	switch v := val.(type) {
	case reflect.Type:
		return reflect.Zero(v)
	case reflect.Value:
		return v
	default:
		return reflect.ValueOf(v)
	}
}

func AnyOf(val interface{}) any {
	switch v := val.(type) {
	case reflect.Type:
		return reflect.Zero(v).Interface()
	case reflect.Value:
		return v.Interface()
	default:
		return v
	}
}

func KindOf(val interface{}) reflect.Kind {
	return ValueOf(val).Kind()
}

func SliceOf(typ interface{}) reflect.Type {
	return reflect.SliceOf(TypeOf(typ))
}

func MakeSlice(val interface{}, nums ...int) any {
	leng := 0
	cap := 0
	if len(nums) > 0 {
		leng = nums[0]
	}
	if len(nums) > 1 {
		cap = nums[1]
	}
	kind := KindOf(val)
	if kind == reflect.Slice {
		return reflect.MakeSlice(TypeOf(val), leng, cap)
	}
	if kind == reflect.Array {
		return reflect.MakeSlice(reflect.SliceOf(TypeOf(val).Elem()), leng, cap)
	}
	return reflect.MakeSlice(SliceOf(val), leng, cap)
}

func MakeSliceOf(value interface{}, nums ...int) any {
	leng := 0
	cap := 0
	if len(nums) > 0 {
		leng = nums[0]
	}
	if len(nums) > 1 {
		cap = nums[1]
	}
	return MakeSlice(SliceOf(value), leng, cap)
}

func GetSliceAt(slice any, index int) any {
	return ValueOf(slice).Index(index).Interface()
}

func SetSliceAt(slice any, index int, newVal any) {
	ValueOf(slice).Index(index).Set(ValueOf(newVal))
}

func Len(val interface{}) int {
	if val == nil {
		return 0
	}
	return ValueOf(val).Len()
}

func Append(s interface{}, x ...interface{}) any {
	for _, v := range x {
		s = reflect.Append(ValueOf(s), ValueOf(v))
	}
	return AnyOf(s)
}

func AppendSlice(s interface{}, x interface{}) any {
	for i := 0; i < Len(x); i++ {
		s = reflect.Append(ValueOf(s), ValueOf(GetSliceAt(x, i)))
	}
	return AnyOf(s)
}

func ElemValue(e interface{}) reflect.Value {
	return ValueOf(e).Elem()
}

func Elem(e interface{}) any {
	return ValueOf(e).Elem().Interface()
}

func SetField(structPointer interface{}, name string, value interface{}) {
	ElemValue(structPointer).FieldByName(name).Set(ValueOf(value))
}

func FieldByName(s interface{}, name string) any {
	return ValueOf(s).FieldByName(name).Interface()
}

func MethodByName(s interface{}, name string) any {
	return ValueOf(s).MethodByName(name).Interface()
}

func MapIndex(m interface{}, key interface{}) any {
	return ValueOf(m).MapIndex(ValueOf(key)).Interface()
}

func PropertyByName(s interface{}, name string) any {
	var zero reflect.Value
	prop := MethodByName(s, name)
	if prop != zero {
		return prop
	}
	prop = FieldByName(s, name)
	if prop != zero {
		return prop
	}
	return MapIndex(s, name)
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	val := ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	case reflect.Interface:
		return val.IsNil() || isNil(val.Elem())
	case reflect.Invalid:
		return true
	default:
		return false
	}
}

func Call(fn interface{}, argss ...[]interface{}) []any {
	var args []interface{}
	if len(args) > 0 {
		args = argss[0]
	}
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = ValueOf(args[i])
	}
	outputs := ValueOf(fn).Call(inputs)
	returns := make([]any, len(outputs))
	for i := range outputs {
		returns[i] = AnyOf(outputs[i])
	}
	return returns
}

func CallMethod(obj interface{}, name string, argss ...[]interface{}) []any {
	return Call(MethodByName(obj, name), argss...)
}

type testStruct struct {
	Name  string
	Count int
	X     []string
}

func ptr[T any](t T) *T {
	return &t
}
func ap[T any](s *[]T, t T) {
	s = ptr(append(*s, t))
}

func MapGetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue ...V) V {
	if m != nil {
		if v, ok := m[key]; ok && any(v) != nil {
			return v
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	var zero V
	return zero
}
