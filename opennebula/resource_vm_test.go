package opennebula

import (
	"testing"
)

func TestGetBaseAttributeFromKeyValue(t *testing.T) {
	cases := []struct {
		key    string
		value  interface{}
		result string
	}{
		{
			key:    "key0",
			value:  12,
			result: "KEY0 = \"12\"",
		},
		{
			key:    "key1",
			value:  "mystring",
			result: "KEY1 = \"mystring\"",
		},
		{
			key:    "empty",
			value:  "",
			result: "EMPTY = \"\"",
		},
		{
			key:    "zero",
			value:  0,
			result: "ZERO = \"0\"",
		},
		{
			key:    "negative",
			value:  -54,
			result: "NEGATIVE = \"-54\"",
		},
	}
	for i, tc := range cases {
		actual := getBaseAttributeFromKeyValue(tc.key, tc.value)
		if actual != tc.result {
			t.Fatalf("Case %d: Expected '%s', got '%s'", i, tc.result, actual)
		}
	}
}
