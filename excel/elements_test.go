package excel

import (
	"reflect"
	"testing"
)

func TestElementsGet(t *testing.T) {
	element1 := make(map[string]string)
	element2 := make(map[string]string)
	element3 := make(map[string]string)

	element1["attr1"] = "value1_1"
	element1["attr2"] = "value1_2"
	element2["attr1"] = "value2_1"
	element3["attr1"] = "value3_1"
	element3["attr3"] = "value3_2"

	elements := Elements([]Attributes{element1, element2, element3})

	if got := elements.Get("dummy", "dummy"); got != nil {
		t.Errorf("Expected element [%v], got [%v]", nil, got)
	}

	if got := elements.Get("attr1", "value1_1"); !reflect.DeepEqual(got, elements[0]) {
		t.Errorf("Expected element [%p], got [%p]", elements[0], got)
	}
	if got := elements.Get("attr1", "value2_1"); !reflect.DeepEqual(got, elements[1]) {
		t.Errorf("Expected element [%v], got [%v]", elements[1], got)
	}
	if got := elements.Get("attr3", "value3_2"); !reflect.DeepEqual(got, elements[2]) {
		t.Errorf("Expected element [%v], got [%v]", elements[2], got)
	}
}
