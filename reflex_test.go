package reflex

import (
	"fmt"
	"testing"
)

func Test_Reflex(t *testing.T) {
	s := []int{1, 2, 3}
	SetSliceAt(s, 1, 42)
	fmt.Println(s) // should print [1, 42, 3]

	ts := &testStruct{
		Name:  "chees",
		Count: 2,
		X:     []string{"poop"},
	}
	SetField(ts, "Name", "asdf")
	SetField(ts, "Count", 22)
	ap(&ts.X, "asdf")
	fmt.Println(ts)

}
