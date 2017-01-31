package binary_pack

import (
	"testing"
	"reflect"
)

func TestBinaryPack_CalcSize(t *testing.T) {
	cases := []struct {
		in []string
		want int
		e bool
	}{
		{[]string{}, 0, false},
		{[]string{"I", "I", "I", "4s"}, 16, false},
		{[]string{"H", "H", "I", "H", "8s", "H"}, 20, false},
		{[]string{"i", "?", "H", "f", "d", "h", "I", "5s"}, 30, false},
		{[]string{"?", "h", "H", "i", "I", "l", "L", "q", "Q", "f", "d", "1s"}, 50, false},
		// Unknown tokens
		{[]string{"a", "b", "c"}, 0, true},
	}

	for _, c := range cases {
		got, err := new(BinaryPack).CalcSize(c.in)

		if err != nil && !c.e {
			t.Errorf("CalcSize(%v) raised %v", c.in, err)
		}

		if err == nil && got != c.want {
			t.Errorf("CalcSize(%v) == %d want %d", c.in, got, c.want)
		}
	}
}

func TestBinaryPack_Pack(t *testing.T) {
	cases := []struct {
		f []string
		a []interface{}
		want []byte
		e bool
	}{
		{[]string{"?", "?"}, []interface{}{true, false}, []byte{1, 0}, false},
		{[]string{"h", "h", "h"}, []interface{}{0, 5, -5},
			[]byte{0, 0, 0, 5, 255, 251}, false},
		{[]string{"H", "H", "H"}, []interface{}{0, 5, 2300},
			[]byte{0, 0, 0, 5, 8, 252}, false},
		{[]string{"i", "i", "i"}, []interface{}{0, 5, -5},
			[]byte{0, 0, 0, 0, 0, 0, 0, 5, 255, 255, 255, 251}, false},
		{[]string{"I", "I", "I"}, []interface{}{0, 5, 2300},
			[]byte{0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 8, 252}, false},
		{[]string{"f", "f", "f"}, []interface{}{float32(0.0), float32(5.3), float32(-5.3)},
			[]byte{0, 0, 0, 0, 64, 169, 153, 154, 192, 169, 153, 154}, false},
		{[]string{"d", "d", "d"}, []interface{}{0.0, 5.3, -5.3},
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 64, 21, 51, 51, 51, 51, 51, 51, 192, 21, 51, 51, 51, 51, 51, 51}, false},
		{[]string{"1s", "2s", "10s"}, []interface{}{"a", "bb", "1234567890"},
			[]byte{97, 98, 98, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48}, false},
		{[]string{"I", "I", "I", "4s"}, []interface{}{1, 2, 4, "DUMP"},
			[]byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 4, 68, 85, 77, 80}, false},
		// Wrong format length
		{[]string{"I", "I", "I", "4s"}, []interface{}{1, 4, "DUMP"}, nil, true},
		// Wrong format token
		{[]string{"I", "a", "I", "4s"}, []interface{}{1, 2, 4, "DUMP"}, nil, true},
	}

	for _, c := range cases {
		got, err := new(BinaryPack).Pack(c.f, c.a)

		if err != nil && !c.e {
			t.Errorf("Pack(%v, %v) raised %v", c.f, c.a, err)
		}

		if err == nil && !reflect.DeepEqual(got, c.want) {
			t.Errorf("Pack(%v, %v) == %v want %v", c.f, c.a, got, c.want)
		}
	}
}

func TestBinaryPack_UnPack(t *testing.T) {
	cases := []struct {
		f []string
		a []byte
		want []interface{}
		e bool
	}{
		{[]string{"?", "?"}, []byte{1, 0}, []interface{}{true, false}, false},
		{[]string{"h", "h", "h"}, []byte{0, 0, 0, 5, 255, 251},
			[]interface{}{0, 5, -5}, false},
		{[]string{"H", "H", "H"}, []byte{0, 0, 0, 5, 8, 252},
			[]interface{}{0, 5, 2300}, false},
		{[]string{"i", "i", "i"}, []byte{0, 0, 0, 0, 0, 0, 0, 5, 255, 255, 255, 251},
			[]interface{}{0, 5, -5}, false},
		{[]string{"I", "I", "I"}, []byte{0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 8, 252},
			[]interface{}{0, 5, 2300}, false},
		{[]string{"f", "f", "f"},
			[]byte{0, 0, 0, 0, 64, 169, 153, 154, 192, 169, 153, 154},
			[]interface{}{float32(0.0), float32(5.3), float32(-5.3)}, false},
		{[]string{"d", "d", "d"},
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 64, 21, 51, 51, 51, 51, 51, 51, 192, 21, 51, 51, 51, 51, 51, 51},
			[]interface{}{0.0, 5.3, -5.3}, false},
		{[]string{"1s", "2s", "10s"},
			[]byte{97, 98, 98, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			[]interface{}{"a", "bb", "1234567890"}, false},
		{[]string{"I", "I", "I", "4s"},
			[]byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 4, 68, 85, 77, 80},
			[]interface{}{1, 2, 4, "DUMP"}, false},
		// Wrong format length
		{[]string{"I", "I", "I", "4s", "H"}, []byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80},
			nil, true},
		// Wrong format token
		{[]string{"I", "a", "I", "4s"}, []byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80},
			nil, true},
	}

	for _, c := range cases {
		got, err := new(BinaryPack).UnPack(c.f, c.a)

		if err != nil && !c.e {
			t.Errorf("UnPack(%v, %v) raised %v", c.f, c.a, err)
		}

		if err == nil && !reflect.DeepEqual(got, c.want) {
			t.Errorf("UnPack(%v, %v) == %v want %v", c.f, c.a, got, c.want)
		}
	}
}
