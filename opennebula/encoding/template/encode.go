package template

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// Marshal returns the TEMPLATE encoding of v.
func Marshal(d interface{}) (string, error) {
	v := reflect.ValueOf(d).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("%s not supported", t.Kind().String())
	}

	attributes := make([]string, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := v.Type().Field(i)
		if isTypeAvailable(vField) && !isEmptyValue(vField) {
			attributes = append(attributes, typeField(tField)+"="+reflectValue(vField))
		}
	}

	return strings.Join(attributes, ", "), nil
}

// MarshalWraps returns the TEMPLATE encoding of v with Wrap
func MarshalWrap(d interface{}) (string, error) {
	name := strings.ToUpper(reflect.ValueOf(d).Elem().Type().Name())
	tpl, err := Marshal(d)
	if err != nil {
		return tpl, err
	}
	return decorate(name, tpl), nil
}

func decorate(name string, attrs string) string {
	response := ""

	if name != "" {
		response += name + " = ["
	}

	response += attrs

	if name != "" {
		response += "]"
	}
	return response
}

func isTypeAvailable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int,
		reflect.String:
		return true
	default:
		return false
	}

}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

func reflectValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int:
		return strconv.Itoa(v.Interface().(int))
	case reflect.String:
		return v.Interface().(string)
	}
	return ""
}

func typeField(f reflect.StructField) string {
	val := f.Name
	tag := f.Tag.Get("template")
	if isValidTag(tag) {
		val = tag
	}
	return strings.ToUpper(val)
}
