package editor

import "testing"

func TestIsSymbol(t *testing.T) {
	testSlice := []byte{
		'a',
		'-',
		'\t',
		'x',
	}

	want := []bool{false, true, true, false}

	for i := range testSlice {
		got := isSymbol(testSlice[i])
		if got != want[i] {
			t.Errorf("row %d: got %t, want %t", i, got, want[i])
		}
	}
}

func TestRowContainsLetterOrDigit(t *testing.T) {
	testEditor := Editor{
		row: []EditorRow{
			{chars: []byte{'a', 'b', 'c'}},                // true
			{chars: []byte{}},                             // false
			{chars: []byte{'-'}},                          // false
			{chars: []byte{'f', 'o', 'o'}},                // true
			{chars: []byte{' ', ' ', '-', '/', ' ', 'a'}}, // true
		},
	}

	want := []bool{true, false, false, true, true}

	for i := range testEditor.row {
		got := rowContainsLetterOrDigit(&testEditor.row[i])
		if got != want[i] {
			t.Errorf("row %d: got %t, want %t", i, got, want[i])
		}
	}
}
