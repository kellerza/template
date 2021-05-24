package template

import "testing"

func TestExec(t *testing.T) {
	tm := New("tst")
	tm.T.Parse("{{ .a }}")

	v := map[string]interface{}{
		"a": "a",
	}

	res, err := tm.Execute(v)
	if err != nil {
		t.Error(err)
	}
	if string(res) != "a" {
		t.Errorf("no expected: %s", res)
	}

}

func TestExecT(t *testing.T) {
	tm := New("tst")
	tm.SearchPath = []string{"./test_data"}

	v := map[string]interface{}{
		"a": "a",
	}

	res, err := tm.ExecuteTemplate("tst.tmpl", v)
	if err != nil {
		t.Error(err)
	}
	if string(res) != "begin\na\nend" {
		t.Errorf("not expected: %s", res)
	}

}
