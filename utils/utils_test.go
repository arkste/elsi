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
