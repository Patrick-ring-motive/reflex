package reflex

import (
	"reflect"
	"sync"
	"testing"
)

type sample struct {
	Name   string
	Count  int
	hidden string
}

func (sample) Greet(prefix string) string {
	return prefix + "world"
}

func (sample) Multiply(value int) int {
	return value * 2
}

type stringAlias string

func requirePanic(t *testing.T, action func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	action()
}

func TestTypeValueAnyAndKind(t *testing.T) {
	if TypeOf(1) != reflect.TypeOf(1) {
		t.Fatal("TypeOf should return dynamic type")
	}
	if TypeOf(reflect.TypeOf("")) != reflect.TypeOf("") {
		t.Fatal("TypeOf should retain reflect.Type")
	}
	if TypeOf(reflect.ValueOf(true)) != reflect.TypeOf(true) {
		t.Fatal("TypeOf should unwrap reflect.Value")
	}
	if TypeOf(reflect.Value{}) != nil {
		t.Fatal("invalid reflect.Value should have no type")
	}
	var nilType reflect.Type
	if TypeOf(nilType) != nil {
		t.Fatal("nil reflect.Type should remain nil")
	}

	if ValueOf(reflect.TypeOf(1)).Int() != 0 {
		t.Fatal("ValueOf should create zero value for reflect.Type")
	}
	if ValueOf(nilType).IsValid() {
		t.Fatal("nil reflect.Type should create an invalid value")
	}
	if ValueOf(reflect.ValueOf("value")).String() != "value" {
		t.Fatal("ValueOf should retain reflect.Value")
	}
	if ValueOf("value").String() != "value" {
		t.Fatal("ValueOf should reflect ordinary values")
	}

	if AnyOf(reflect.TypeOf(1)) != 0 {
		t.Fatal("AnyOf should unwrap zero value for reflect.Type")
	}
	if AnyOf(nilType) != nil || AnyOf(reflect.Value{}) != nil {
		t.Fatal("AnyOf should return nil for invalid reflection inputs")
	}
	if AnyOf(reflect.ValueOf("value")) != "value" || AnyOf("value") != "value" {
		t.Fatal("AnyOf should preserve ordinary values")
	}
	if KindOf(nil) != reflect.Invalid || KindOf(1) != reflect.Int {
		t.Fatal("KindOf returned unexpected kind")
	}
}

func TestSliceHelpers(t *testing.T) {
	if SliceOf(1) != reflect.TypeOf([]int{}) || SliceOf(reflect.Value{}) != nil {
		t.Fatal("SliceOf returned unexpected type")
	}

	fromSlice := MakeSlice([]int{}, 2).([]int)
	if len(fromSlice) != 2 || cap(fromSlice) != 2 {
		t.Fatalf("MakeSlice slice result = len %d cap %d", len(fromSlice), cap(fromSlice))
	}
	fromArray := MakeSlice([2]string{}, 1, 3).([]string)
	if len(fromArray) != 1 || cap(fromArray) != 3 {
		t.Fatalf("MakeSlice array result = len %d cap %d", len(fromArray), cap(fromArray))
	}
	fromValue := MakeSlice(0, 1, 2).([]int)
	if len(fromValue) != 1 || cap(fromValue) != 2 {
		t.Fatalf("MakeSlice value result = len %d cap %d", len(fromValue), cap(fromValue))
	}
	if result := MakeSliceOf(0).([]int); len(result) != 0 || cap(result) != 0 {
		t.Fatal("MakeSliceOf returned unexpected slice")
	}

	requirePanic(t, func() { MakeSlice(nil) })
	requirePanic(t, func() { MakeSlice(0, 1, 1, 1) })
}

func TestCollectionMutationAndElements(t *testing.T) {
	values := []int{1, 2}
	SetSliceAt(values, 1, int8(4))
	if GetSliceAt(values, 1) != 4 {
		t.Fatal("SetSliceAt should convert assignable values")
	}

	pointer := 1
	pointers := []*int{&pointer}
	SetSliceAt(pointers, 0, nil)
	if pointers[0] != nil {
		t.Fatal("SetSliceAt should assign nil to nilable targets")
	}
	requirePanic(t, func() { SetSliceAt(values, 0, "invalid") })
	requirePanic(t, func() { SetSliceAt(values, 0, nil) })

	if Len(nil) != 0 || Len(reflect.Value{}) != 0 || Len("abc") != 3 {
		t.Fatal("Len returned unexpected result")
	}
	appended := Append([]int{1}, int8(2)).([]int)
	if !reflect.DeepEqual(appended, []int{1, 2}) {
		t.Fatalf("Append result = %v", appended)
	}
	if result := Append([]int{1}).([]int); !reflect.DeepEqual(result, []int{1}) {
		t.Fatalf("Append without values = %v", result)
	}
	appendedSlice := AppendSlice([]int{1}, []int8{2, 3}).([]int)
	if !reflect.DeepEqual(appendedSlice, []int{1, 2, 3}) {
		t.Fatalf("AppendSlice result = %v", appendedSlice)
	}

	value := 5
	if ElemValue(&value).Int() != 5 || Elem(&value) != 5 {
		t.Fatal("element helpers returned unexpected value")
	}

	subject := sample{Name: "before", Count: 1}
	SetField(&subject, "Count", int8(7))
	if subject.Count != 7 {
		t.Fatal("SetField should convert assignable values")
	}
}

func TestLookupAndNilHelpers(t *testing.T) {
	subject := sample{Name: "field", hidden: "private"}
	if FieldByName(subject, "Name") != "field" || FieldByName(&subject, "Name") != "field" {
		t.Fatal("FieldByName should find exported fields")
	}
	if FieldByName(subject, "hidden") != nil || FieldByName(subject, "missing") != nil {
		t.Fatal("FieldByName should hide inaccessible and missing fields")
	}
	var nilSubject *sample
	if FieldByName(nilSubject, "Name") != nil || FieldByName(1, "Name") != nil {
		t.Fatal("FieldByName should reject invalid inputs")
	}

	method, ok := MethodByName(subject, "Greet").(func(string) string)
	if !ok || method("hello ") != "hello world" {
		t.Fatal("MethodByName should find methods")
	}
	if MethodByName(subject, "missing") != nil || MethodByName(nil, "missing") != nil {
		t.Fatal("MethodByName should return nil for missing methods")
	}

	values := map[string]int{"found": 9}
	if MapIndex(values, stringAlias("found")) != 9 {
		t.Fatal("MapIndex should convert keys")
	}
	if MapIndex(values, "missing") != nil || MapIndex(values, nil) != nil || MapIndex(1, "found") != nil {
		t.Fatal("MapIndex should return nil for missing and invalid lookups")
	}

	if PropertyByName(subject, "Greet") == nil || PropertyByName(subject, "Name") != "field" || PropertyByName(values, "found") != 9 || PropertyByName(subject, "missing") != nil {
		t.Fatal("PropertyByName returned unexpected value")
	}

	var nilPointer *int
	if !IsNil(nil) || !IsNil(nilPointer) || !IsNil(reflect.Value{}) {
		t.Fatal("IsNil should find direct nil values")
	}
	var dynamic any = nilPointer
	wrapped := reflect.ValueOf(&dynamic).Elem()
	if !IsNil(wrapped) || IsNil(0) {
		t.Fatal("IsNil should handle wrapped interfaces and values")
	}
}

func TestCallCloseAndMapDefault(t *testing.T) {
	noArguments := func() int { return 1 }
	if result := Call(noArguments); len(result) != 1 || result[0] != 1 {
		t.Fatalf("Call no-argument result = %v", result)
	}
	add := func(left, right int) int { return left + right }
	if result := Call(add, []any{int8(2)}, []any{int8(3)}); len(result) != 1 || result[0] != 5 {
		t.Fatalf("Call converted split arguments result = %v", result)
	}
	if result := Call(add, []any{int8(2)}); len(result) != 1 || result[0] != 2 {
		t.Fatalf("Call missing argument result = %v", result)
	}
	if result := Call(add, []any{"invalid", int8(3), 99}); len(result) != 1 || result[0] != 3 {
		t.Fatalf("Call incompatible or extra argument result = %v", result)
	}
	sum := func(start int, values ...int) int {
		for _, value := range values {
			start += value
		}
		return start
	}
	if result := Call(sum, []any{1}, []any{int8(2), int8(3)}); len(result) != 1 || result[0] != 6 {
		t.Fatalf("Call variadic result = %v", result)
	}
	if result := CallMethod(sample{}, "Multiply", []any{int8(4)}); len(result) != 1 || result[0] != 8 {
		t.Fatalf("CallMethod result = %v", result)
	}

	values := map[string]int{"answer": 42}
	if result := Call(values, []any{stringAlias("answer")}); len(result) != 1 || result[0] != 42 {
		t.Fatalf("Call map getter result = %v", result)
	}
	if result := Call(values); len(result) != 1 || result[0] != nil {
		t.Fatalf("Call map getter without key result = %v", result)
	}

	slice := []string{"zero", "one"}
	array := [2]int{7, 8}
	if result := Call(slice, []any{int8(1)}); len(result) != 1 || result[0] != "one" {
		t.Fatalf("Call slice getter result = %v", result)
	}
	if result := Call(array, []any{0}); len(result) != 1 || result[0] != 7 {
		t.Fatalf("Call array getter result = %v", result)
	}
	if result := Call(slice); len(result) != 1 || result[0] != nil {
		t.Fatalf("Call slice getter without index result = %v", result)
	}
	if result := Call(slice, []any{"invalid"}); len(result) != 1 || result[0] != nil {
		t.Fatalf("Call slice with invalid index result = %v", result)
	}
	if result := Call(slice, []any{-1}); len(result) != 1 || result[0] != nil {
		t.Fatalf("Call slice with out-of-range index result = %v", result)
	}

	var store sync.Map
	store.Store("answer", 42)
	if result := Call(&store, []any{"answer"}); len(result) != 2 || result[0] != 42 || result[1] != true {
		t.Fatalf("Call sync.Map getter result = %v", result)
	}
	if result := Call(&store); len(result) != 2 || result[0] != nil || result[1] != false {
		t.Fatalf("Call sync.Map getter without key result = %v", result)
	}
	var nilStore *sync.Map
	if result := Call(nilStore, []any{"answer"}); len(result) != 2 || result[0] != nil || result[1] != false {
		t.Fatalf("Call nil sync.Map getter result = %v", result)
	}

	var nilFunction func()
	if Call(nil) != nil || Call(1) != nil || Call(nilFunction) != nil || CallMethod(sample{}, "missing") != nil {
		t.Fatal("Call should return nil for invalid functions")
	}
	if Call(func() { panic("failure") }) != nil {
		t.Fatal("Call should recover invocation panics")
	}

	channel := make(chan int)
	Close(channel)
	if _, open := <-channel; open {
		t.Fatal("Close should close channel")
	}

	if MapGetOrDefault(map[string]int{"found": 3}, "found", 8) != 3 {
		t.Fatal("MapGetOrDefault should return found value")
	}
	if MapGetOrDefault(map[string]int(nil), "missing", 8) != 8 {
		t.Fatal("MapGetOrDefault should return default for nil map")
	}
	var nilPointer *int
	fallback := new(int)
	if MapGetOrDefault(map[string]*int{"nil": nilPointer}, "nil", fallback) != fallback {
		t.Fatal("MapGetOrDefault should return default for typed nil value")
	}
	if MapGetOrDefault(map[string]int{}, "missing") != 0 {
		t.Fatal("MapGetOrDefault should return zero without a default")
	}
}
