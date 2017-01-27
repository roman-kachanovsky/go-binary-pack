package binary_pack

import "testing"

func TestBinaryPack_CalcSize(t *testing.T) {
	cases := []struct {
		in []string
		want int
	}{
		{[]string{}, 0},
		{[]string{"I", "I", "I", "4s"}, 16},
		{[]string{"H", "H", "I", "H", "8s", "H"}, 20},
		{[]string{"i", "q", "?", "H", "f", "d", "h", "I", "5s"}, 30},
	}

	for _, c := range cases {
		got := new(BinaryPack).CalcSize(c.in)

		if got != c.want {
			t.Errorf("CalcSize(%v) == %d want %d", c.in, got, c.want)
		}
	}
}
