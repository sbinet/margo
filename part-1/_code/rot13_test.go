// STARTIMPORT OMIT

package rot13 // HLxxx

import (
	"reflect"
	"testing" // HLxxx
)

type testCase struct {
	text []byte
	want []byte
}

// ENDIMPORT OMIT

// STARTTESTCASES OMIT

var cases = []testCase{
	{
		text: []byte("Lbh penpxrq gur pbqr!"),
		want: []byte("You cracked the code!"),
	},
	{
		text: []byte("hello"),
		want: []byte("uryyb"),
	},
	{
		text: []byte("abcdefghijklmnopqrstuvwxyz"),
		want: []byte("nopqrstuvwxyzabcdefghijklm"),
	},
	{
		text: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		want: []byte("NOPQRSTUVWXYZABCDEFGHIJKLM"),
	},
}

// ENDTESTCASES OMIT

func TestRot13(t *testing.T) {
	for _, table := range cases {
		o := make([]byte, len(table.text))
		for i, b := range table.text {
			o[i] = rot13(b) // HLxxx
		}

		if !reflect.DeepEqual(o, table.want) {
			t.Fatalf("invalid rot13\ngot =%q\nwant=%q\n", string(o), string(table.want))
		}
	}
}
