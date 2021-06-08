package template

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// The following functions will be available when rendering the template.
//
// Function slice & index overwrites the standard Go functions, but is compatible to the standard library
var Funcs = map[string]interface{}{
	"optional": Optional,
	"expect":   Expect,
	"ip":       Ip,
	"ipmask":   Ipmask,
	"default":  Default,
	"contains": Contains,
	"index":    Index,
	"join":     Join,
	"slice":    Slice,
	"split":    Split,
}

// Get an int from a relfect.Value and if this was a valid int
func parseInt(index reflect.Value) (int, bool) {
	switch index.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(index.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int(index.Uint()), true
	default:
		return 0, false
	}
}

// Test if a string contains a substring
func Contains(substr, str string) (interface{}, error) {
	return strings.Contains(fmt.Sprintf("%v", str), fmt.Sprintf("%v", substr)), nil
}

// Return a default value if a value is not available
func Default(def, val interface{}) (interface{}, error) {

	switch v := val.(type) {
	case nil:
		return def, nil
	case string:
		if v == "" {
			return def, nil
		}
	case bool:
		if !v {
			return def, nil
		}
	}
	// if val == nil {
	// 	return def, nil
	// }

	// If we have a input value, do some type checking
	tval, tdef := typeof(val), typeof(def)
	if tval == "string" && tdef == "int" {
		if _, err := strconv.Atoi(val.(string)); err == nil {
			tval = "int"
		}
		if tdef == "str" {
			if _, err := strconv.Atoi(def.(string)); err == nil {
				tdef = "int"
			}
		}
	}
	if tdef != tval {
		return val, fmt.Errorf("expected type %v, got %v (value=%v)", tdef, tval, val)
	}

	// Return the value
	return val, nil
}

// The indexes can either follow the value, or be before the value (supporting pipe)
// Negative indexes are allowed and will be the offset from the length
func Index(args ...reflect.Value) (reflect.Value, error) {
	if len(args) < 2 {
		return reflect.Value{}, fmt.Errorf("at least 2 parameters expected")
	}
	indexes := make([]reflect.Value, len(args)-1)

	// idx=0: support the built-in parameter order
	// idx=1: support parameter order with value last (to pipe)
	var item reflect.Value
	for offs := 0; offs < 2; offs++ {
		switch offs {
		case 0:
			item = indirectInterface(args[(len(args) - 1)])
		case 1:
			item = indirectInterface(args[0])
		}
		switch item.Kind() {
		case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
			for i := 0; i < len(args)-1; i++ {
				indexes[i] = args[i+offs]
			}
			return index_builtin(item, indexes...)
		}
	}

	return reflect.Value{}, fmt.Errorf("expected an array, slice, string or map and an index %s %s", args[0].Kind(), args[len(args)-1].Kind())
}

// Return only the IP address from a value containing an IP/mask
func Ip(val interface{}) (interface{}, error) {
	s := fmt.Sprintf("%v", val)
	a := strings.Split(s, "/")
	return a[0], nil
}

// Return only the mask from a value containing an IP/mask
func Ipmask(val interface{}) (interface{}, error) {
	s := fmt.Sprintf("%v", val)
	a := strings.Split(s, "/")
	return a[1], nil
}

// Joins an array of values or slice using the specified separator
func Join(sep string, val reflect.Value) (interface{}, error) {
	if sep == "" {
		sep = " "
	}
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() == 0 {
			return "", nil
		}
		var s strings.Builder
		i := 0
		for ; i < val.Len()-1; i++ {
			fmt.Fprintf(&s, "%v", val.Index(i))
			fmt.Fprint(&s, sep)
		}
		fmt.Fprintf(&s, "%v", val.Index(i))
		return s.String(), nil
	}
	return nil, fmt.Errorf("expected array [], got %v [%s]", val, val.Kind())
}

// Slicing.
//
// slice returns the result of text/template's [slice](https://golang.org/pkg/text/template/#hdr-Functions)
// if that fails, it attempts an alternative implementation, the the first 2 parameters
// are indexes followed by the value.
// Negative indexes are allowed and will be the offset from the length
func Slice(item reflect.Value, indexes ...reflect.Value) (reflect.Value, error) {
	// call the internal function
	res, err := slice_builtin(item, indexes...)
	if err == nil {
		return res, nil
	}
	if len(indexes) != 2 {
		return reflect.Value{}, err
	}

	// accept the value as the last argument to support pipes
	start, ok1 := parseInt(item)
	end, ok2 := parseInt(indexes[0])
	if !ok1 || !ok2 {
		return reflect.Value{}, err
	}
	val := indirectInterface(indexes[1])
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice:
		if start < 0 {
			start += val.Len()
		}
		if end <= 0 {
			end += val.Len()
		}
		return val.Slice(start, end), nil
	}
	return reflect.Value{}, fmt.Errorf("not an array, string or slice")
}

// Split a string using the separator
func Split(sep string, val interface{}) (interface{}, error) {
	// Start and end values
	if val == nil {
		return []interface{}{}, nil
	}
	if sep == "" {
		sep = " "
	}

	v := fmt.Sprintf("%v", val)

	res := strings.Split(v, sep)
	r := make([]interface{}, len(res))
	for i, p := range res {
		r[i] = p
	}
	return r, nil
}
