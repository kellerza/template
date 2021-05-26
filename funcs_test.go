package template

import (
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestRender1(t *testing.T) {

	var test_set = map[string][]string{
		// empty values
		"default 0 .x":      {"0"},
		".x | default 0":    {"0"},
		"default 0 \"\"":    {"0"},
		"default 0 false":   {"0"},
		"false | default 0": {"0"},
		// ints pass through ok
		"default 1 0": {"0"},
		// errors
		"default .x":     {"", "wrong number"},
		"default .x 1 1": {"", "wrong number"},
		// type check
		"default 0 .i5":   {"5"},
		"default 0 .sA":   {"", "expected type int"},
		`default "5" .sA`: {"A"},

		`contains "." .sAAA`:   {"true"},
		`.sAAA | contains "."`: {"true"},
		`contains "." .sA`:     {"false"},
		`.sA | contains "."`:   {"false"},

		`split "." "a.a"`:  {"[a a]"},
		`split " " "a bb"`: {"[a bb]"},

		`ip "1.1.1.1/32"`:     {"1.1.1.1"},
		`"1.1.1.1" | ip`:      {"1.1.1.1"},
		`ipmask "1.1.1.1/32"`: {"32"},
		`"1.1.1.1/32" | split "/" | slice 0 1 | join ""`:  {"1.1.1.1"},
		`"1.1.1.1/32" | split "/" | slice 1 2 | join ""`:  {"32"},
		`"1.1.1.1/32" | split "/" | slice -1 0 | join ""`: {"32"},

		`slice "abc" 1 2`:  {"b"}, // built in parameter order
		`slice  1 2 "abc"`: {"b"}, // parametr order for pipe mode

		`split " " "a bb" | join "-"`: {"a-bb"},
		`split "" ""`:                 {"[]"},
		`split "abc" ""`:              {"[]"},

		`index .iparr 1`: {"32"},               // built in parameter order
		`index .iparr 2`: {"", "out of range"}, // built in parameter order
		`index 1 .iparr`: {"32"},               // parameter order for pipe mode
		`index 1 1`:      {"", "expected"},     // expected array
		`index 1`:        {"", "at least"},     // too few args

		`"1.1.1.1/32" | split "/" | index 1`:  {"32"},
		`"1.1.1.1/32" | split "/" | index -1`: {"32"},
		`"1.1.1.1/32" | split "/" | index -2`: {"1.1.1.1"},
		`"1.1.1.1/32" | split "/" | index -3`: {"", "out of range"},
		`"1.1.1.1/32" | split "/" | index 2`:  {"", "out of range"},
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

func test_check(templateS string, vars map[string]interface{}, exp ...string) error {
	var buf strings.Builder
	ts := fmt.Sprintf("{{ %v }}", strings.Trim(templateS, "{} "))
	tem, err0 := template.New("").Funcs(Funcs).Parse(ts)
	if err0 != nil {
		return fmt.Errorf("invalid template")
	}

	template_err := tem.Execute(&buf, vars)

	res := buf.String()

	e := []string{fmt.Sprintf(`%v = "%v", error=%v`, ts, res, template_err)}

	// Check value
	if res != exp[0] {
		e = append(e, fmt.Sprintf("- expected value = %v", exp[0]))
	}

	// Check errors
	if len(exp) > 1 {
		ee := fmt.Sprintf("- expected error with %s", exp[1])
		if template_err == nil {
			e = append(e, ee)
		} else if !strings.Contains(template_err.Error(), exp[1]) {
			e = append(e, ee)
		}
	} else if template_err != nil {
		e = append(e, "- no error expected")
	}
	if len(e) > 1 {
		return fmt.Errorf(strings.Join(e, "\n"))
	}
	return nil
}
