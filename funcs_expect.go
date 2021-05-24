package template

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"inet.af/netaddr"
)

func typeof(val interface{}) string {
	switch val.(type) {
	case string:
		return "string"
	case int, int16, int32:
		return "int"
	}
	return ""
}

// interface{} version of parseInt
func parseInt_i(val interface{}) (int, bool) {
	if i, err := strconv.Atoi(fmt.Sprintf("%v", val)); err == nil {
		return i, true
	}
	return 0, false
}

func expectFunc(val interface{}, format string) (interface{}, error) {
	t := typeof(val)
	vals := fmt.Sprintf("%s", val)

	// known formats
	switch format {
	case "str", "string":
		if t == "string" {
			return "", nil
		}
		return "", fmt.Errorf("string expected, got %s (%v)", t, val)
	case "int":
		if _, ok := parseInt_i(val); ok {
			return "", nil
		}
		return "", fmt.Errorf("int expected, got %s (%v)", t, val)
	case "ip":
		if _, err := netaddr.ParseIPPrefix(vals); err == nil {
			return "", nil
		}
		return "", fmt.Errorf("IP/mask expected, got %v", val)
	}

	// try range
	if matched, _ := regexp.MatchString(`\d+-\d+`, format); matched {
		iv, ok := parseInt_i(val)
		if !ok {
			return "", fmt.Errorf("int expected, got %s (%v)", t, val)
		}
		r := strings.Split(format, "-")
		i0, _ := parseInt_i(r[0])
		i1, _ := parseInt_i(r[1])
		if i1 < i0 {
			i0, i1 = i1, i0
		}
		if i0 <= iv && iv <= i1 {
			return "", nil
		}
		return "", fmt.Errorf("value (%d) expected to be in range %d-%d", iv, i0, i1)
	}

	// Try regex
	matched, err := regexp.MatchString(format, vals)
	if err != nil || !matched {
		return "", fmt.Errorf("value %s does not match regex %s %v", vals, format, err)
	}

	return "", nil
}
