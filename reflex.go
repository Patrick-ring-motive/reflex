package reflex

import (
	"reflect"
)

// TypeOf returns value's dynamic type. A reflect.Type is returned unchanged;
// an invalid reflect.Value and an untyped nil return nil.
func TypeOf(value any) reflect.Type {
	switch value := value.(type) {
	case reflect.Type:
		return value
	case reflect.Value:
		if !value.IsValid() {
			return nil
		}
		return value.Type()
	default:
		return reflect.TypeOf(value)
	}
}

// ValueOf returns value as a reflect.Value. A reflect.Type becomes its zero
// value, while nil becomes an invalid reflect.Value.
func ValueOf(value any) reflect.Value {
	switch value := value.(type) {
	case reflect.Type:
		return reflect.Zero(value)
	case reflect.Value:
		return value
	default:
		return reflect.ValueOf(value)
	}
}

// AnyOf unwraps a reflection value or type. Invalid values and nil map to nil.
// Other values are returned unchanged.
func AnyOf(value any) any {
	switch value := value.(type) {
	case reflect.Type:
		return reflect.Zero(value).Interface()
	case reflect.Value:
		if !value.IsValid() {
			return nil
		}
		return value.Interface()
	default:
		return value
	}
}

// KindOf returns value's reflect.Kind. It returns reflect.Invalid for nil and
// invalid reflection values.
func KindOf(value any) reflect.Kind {
	return ValueOf(value).Kind()
}

// SliceOf returns a slice type whose element type is value's type. It returns
// nil when value has no type.
func SliceOf(value any) reflect.Type {
	elementType := TypeOf(value)
	if elementType == nil {
		return nil
	}
	return reflect.SliceOf(elementType)
}

// MakeSlice creates a slice using value's type. If value is a slice, its type
// is retained; if it is an array, its element type is used. Size accepts zero,
// one (length), or two (length and capacity) values.
func MakeSlice(value any, size ...int) any {
	length, capacity := sliceSize(size)
	sliceType := TypeOf(value)
	if sliceType == nil {
		panic("reflex: cannot make a slice from an untyped nil")
	}

	switch sliceType.Kind() {
	case reflect.Array:
		sliceType = reflect.SliceOf(sliceType.Elem())
	case reflect.Slice:
	default:
		sliceType = reflect.SliceOf(sliceType)
	}

	return reflect.MakeSlice(sliceType, length, capacity).Interface()
}

func sliceSize(size []int) (int, int) {
	switch len(size) {
	case 0:
		return 0, 0
	case 1:
		return size[0], size[0]
	case 2:
		return size[0], size[1]
	default:
		panic("reflex: MakeSlice accepts at most length and capacity")
	}
}

// MakeSliceOf creates a slice using the type returned by SliceOf(value).
func MakeSliceOf(value any, size ...int) any {
	return MakeSlice(SliceOf(value), size...)
}

// GetSliceAt returns slice[index]. Reflection panics for an invalid index or
// a value that cannot be indexed.
func GetSliceAt(slice any, index int) any {
	return AnyOf(ValueOf(slice).Index(index))
}

// SetSliceAt assigns value to slice[index]. Assignable and convertible values
// are accepted; invalid assignments panic.
func SetSliceAt(slice any, index int, value any) {
	target := ValueOf(slice).Index(index)
	target.Set(assignableValue(value, target.Type()))
}

func assignableValue(value any, targetType reflect.Type) reflect.Value {
	converted, ok := convertibleValue(value, targetType)
	if !ok {
		panic("reflex: value cannot be assigned to target type")
	}
	return converted
}

func convertibleValue(value any, targetType reflect.Type) (reflect.Value, bool) {
	source := ValueOf(value)
	if !source.IsValid() {
		switch targetType.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			return reflect.Zero(targetType), true
		default:
			return reflect.Value{}, false
		}
	}
	if source.Type().AssignableTo(targetType) {
		return source, true
	}
	if source.Type().ConvertibleTo(targetType) {
		return source.Convert(targetType), true
	}
	return reflect.Value{}, false
}

// Len returns value's length, or zero for nil and invalid reflection values.
// Reflection panics when value does not have a length.
func Len(value any) int {
	reflectionValue := ValueOf(value)
	if !reflectionValue.IsValid() {
		return 0
	}
	return reflectionValue.Len()
}

// Append appends values to a slice and returns a slice of the original type.
// Each value may be assignable or convertible to the slice element type.
func Append(slice any, values ...any) any {
	result := ValueOf(slice)
	for _, value := range values {
		result = reflect.Append(result, assignableValue(value, result.Type().Elem()))
	}
	return AnyOf(result)
}

// AppendSlice appends every indexed value from values to a slice and returns a
// slice of the original type.
func AppendSlice(slice any, values any) any {
	result := ValueOf(slice)
	for index := 0; index < Len(values); index++ {
		result = reflect.Append(result, assignableValue(GetSliceAt(values, index), result.Type().Elem()))
	}
	return AnyOf(result)
}

// ElemValue returns value's reflected element. Reflection panics when value
// cannot be dereferenced.
func ElemValue(value any) reflect.Value {
	return ValueOf(value).Elem()
}

// Elem returns value's element. Reflection panics when value cannot be
// dereferenced.
func Elem(value any) any {
	return AnyOf(ElemValue(value))
}

// SetField sets the named field on a pointer to a struct. Assignable and
// convertible values are accepted. Reflection panics for inaccessible fields.
func SetField(structPointer any, name string, value any) {
	target := ElemValue(structPointer).FieldByName(name)
	target.Set(assignableValue(value, target.Type()))
}

// FieldByName returns an exported named field or nil when no accessible field
// exists. It accepts a struct or a non-nil pointer to one.
func FieldByName(value any, name string) any {
	structValue := ValueOf(value)
	if structValue.IsValid() && structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return nil
		}
		structValue = structValue.Elem()
	}
	if !structValue.IsValid() || structValue.Kind() != reflect.Struct {
		return nil
	}

	field := structValue.FieldByName(name)
	if !field.IsValid() || !field.CanInterface() {
		return nil
	}
	return AnyOf(field)
}

// MethodByName returns a named method or nil when it does not exist.
func MethodByName(value any, name string) any {
	reflectionValue := ValueOf(value)
	if !reflectionValue.IsValid() {
		return nil
	}
	method := reflectionValue.MethodByName(name)
	if !method.IsValid() || !method.CanInterface() {
		return nil
	}
	return AnyOf(method)
}

// MapIndex returns the value stored under key, or nil when value is not a map,
// key cannot be converted to the map key type, or key is absent.
func MapIndex(value any, key any) any {
	mapValue := ValueOf(value)
	if !mapValue.IsValid() || mapValue.Kind() != reflect.Map {
		return nil
	}

	mapKey, ok := convertibleValue(key, mapValue.Type().Key())
	if !ok {
		return nil
	}
	result := mapValue.MapIndex(mapKey)
	if !result.IsValid() {
		return nil
	}
	return AnyOf(result)
}

// PropertyByName returns a named method, field, or map value, in that order.
// It returns nil when no property is found.
func PropertyByName(value any, name string) any {
	if method := MethodByName(value, name); method != nil {
		return method
	}
	if field := FieldByName(value, name); field != nil {
		return field
	}
	return MapIndex(value, name)
}

// IsNil reports whether value is nil, including typed nil values and invalid
// reflection values.
func IsNil(value any) bool {
	if value == nil {
		return true
	}

	reflectionValue := ValueOf(value)
	if !reflectionValue.IsValid() {
		return true
	}
	switch reflectionValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflectionValue.IsNil()
	case reflect.Interface:
		return reflectionValue.IsNil() || IsNil(reflectionValue.Elem())
	default:
		return false
	}
}

// Call invokes fn with the supplied argument slices and returns its results.
// Maps, arrays, and slices act as getters whose first argument is a key or
// index. Values with sync.Map's Load method return Load's value and found flag.
// Functions convert compatible values, fill missing or incompatible arguments
// with zero values, and ignore extra arguments for non-variadic functions. It
// returns nil when fn is invalid or invocation panics.
func Call(fn any, argumentLists ...[]any) (results []any) {
	defer func() {
		if recover() != nil {
			results = nil
		}
	}()

	var arguments []any
	for _, argumentList := range argumentLists {
		arguments = append(arguments, argumentList...)
	}
	if results, handled := callGetter(fn, arguments); handled {
		return results
	}

	function := ValueOf(fn)
	if !function.IsValid() || function.Kind() != reflect.Func || function.IsNil() {
		return nil
	}

	outputs := function.Call(callInputs(function.Type(), arguments))
	results = make([]any, len(outputs))
	for index, output := range outputs {
		results[index] = AnyOf(output)
	}
	return results
}

type keyValueLoader interface {
	Load(key any) (value any, ok bool)
}

func callGetter(value any, arguments []any) ([]any, bool) {
	reflectionValue := ValueOf(value)
	if !reflectionValue.IsValid() {
		return nil, false
	}
	if reflectionValue.CanInterface() {
		if loader, ok := reflectionValue.Interface().(keyValueLoader); ok {
			if len(arguments) == 0 || (reflectionValue.Kind() == reflect.Ptr && reflectionValue.IsNil()) {
				return []any{nil, false}, true
			}
			loadedValue, found := loader.Load(arguments[0])
			return []any{loadedValue, found}, true
		}
	}

	switch reflectionValue.Kind() {
	case reflect.Map:
		if len(arguments) == 0 {
			return []any{nil}, true
		}
		return []any{MapIndex(value, arguments[0])}, true
	case reflect.Array, reflect.Slice:
		if len(arguments) == 0 {
			return []any{nil}, true
		}
		indexValue, ok := convertibleValue(arguments[0], reflect.TypeOf(0))
		if !ok {
			return []any{nil}, true
		}
		index := int(indexValue.Int())
		if index < 0 || index >= reflectionValue.Len() {
			return []any{nil}, true
		}
		return []any{AnyOf(reflectionValue.Index(index))}, true
	default:
		return nil, false
	}
}

func callInputs(functionType reflect.Type, arguments []any) []reflect.Value {
	inputCount := functionType.NumIn()
	if functionType.IsVariadic() && len(arguments) > inputCount-1 {
		inputCount = len(arguments)
	}

	inputs := make([]reflect.Value, inputCount)
	for index := range inputs {
		targetType := callInputType(functionType, index)
		inputs[index] = reflect.Zero(targetType)
		if index < len(arguments) {
			if converted, ok := convertibleValue(arguments[index], targetType); ok {
				inputs[index] = converted
			}
		}
	}
	return inputs
}

func callInputType(functionType reflect.Type, index int) reflect.Type {
	if functionType.IsVariadic() && index >= functionType.NumIn()-1 {
		return functionType.In(functionType.NumIn() - 1).Elem()
	}
	return functionType.In(index)
}

// CallMethod invokes a named method with the supplied argument slices and
// returns its results. It shares Call's non-panicking behavior.
func CallMethod(value any, name string, argumentLists ...[]any) []any {
	return Call(MethodByName(value, name), argumentLists...)
}

// Close closes channel. Reflection panics when channel is invalid or cannot be
// closed.
func Close(channel any) {
	ValueOf(channel).Close()
}

// MapGetOrDefault returns m[key] when it exists and is non-nil. Otherwise it
// returns the first default value, or V's zero value when no default is given.
func MapGetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue ...V) V {
	if value, ok := m[key]; ok && !IsNil(value) {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	var zero V
	return zero
}
