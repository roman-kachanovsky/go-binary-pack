package binary_pack

import "testing"

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
