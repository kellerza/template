package template

import (
	"testing"
)

func TestRenderExpect(t *testing.T) {

	var test_set = map[string][]string{

		`expect "1.1.1.1/32" "ip"`:   {""},
		`expect "1.1.1.1" "ip"`:      {"", "IP/mask"},
		`expect "1" "0-10"`:          {""},
		`expect "1" "10-10"`:         {"", "range"},
		`expect "1.1" "\\d+\\.\\d+"`: {""},
		`expect 11 "\\d"`:            {""},
		`expect 11 "\\d+"`:           {""},
		`expect "abc" "^[a-z]+$"`:    {""},

		`expect 1 "int"`:    {""},
		`expect 1 "str"`:    {"", "string expected"},
		`expect 1 "string"`: {"", "string expected"},
		`expect .i5 "int"`:  {""},
		`expect "5" "int"`:  {""}, // hasInt
		`expect "aa" "int"`: {"", "int expected"},

		`optional 1 "int"`:   {""},
		`optional .x "int"`:  {""},
		`optional .x "str"`:  {""},
		`optional .i5 "str"`: {""}, // corner case, although it hasInt everything is always a string
	}

	vars := map[string]interface{}{
		"i5":    "5",
		"sA":    "A",
		"sAAA":  "aa.",
		"dot":   ".",
		"space": " ",
		"iparr": []string{"1.1.1.1", "32"},
	}

	for tem, exp := range test_set {
		err := test_check(tem, vars, exp...)
		if err != nil {
			t.Error(err)
		}
	}

}
