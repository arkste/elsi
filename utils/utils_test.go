package utils

import "testing"

func TestCreateHashFromString(t *testing.T) {
	tables := []struct {
		input    string
		expected string
	}{
		{"foobar", "3858f62230ac3c915f300c664312c63f"},
		{"123", "202cb962ac59075b964b07152d234b70"},
	}

	for _, table := range tables {
		newHash := CreateHashFromString(table.input)
		if newHash != table.expected {
			t.Errorf("Hash of %s, got: %s, want: %s.", table.input, newHash, table.expected)
		}
	}
}

func TestIsJSONArray(t *testing.T) {
	tables := []struct {
		input    string
		expected string
	}{
		{"{\"foobar\"}", "{\"foobar\"}"},
		{"[{\"foobar\"}]", "{\"data\":[{\"foobar\"}]}"},
	}

	for _, table := range tables {
		json := IsJSONArray(table.input)
		if json != table.expected {
			t.Errorf("JSON %s should become %s, got %s", table.input, table.expected, json)
		}
	}
}

func TestConvertToJSON(t *testing.T) {
	columns := []string{
		"id",
		"foo",
	}
	values := make([]interface{}, 2)
	values[0] = "1"
	values[1] = "bar"

	expected := "{\"foo\":\"bar\",\"id\":\"1\"}"

	id, returnJSON := ConvertToJSON(columns, values)

	if id != "1" || expected != returnJSON {
		t.Errorf("%v & %v should become %s, got %s", columns, values, expected, returnJSON)
	}
}
