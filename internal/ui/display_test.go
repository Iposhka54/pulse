package ui

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		in  uint64
		out string
	}{
		{0, "0 B"},
		{1024, "1.0 KB"},
		{2000, "2.0 KB"},
		{1_048_575, "1024.0 KB"},
		{1_048_576, "1.0 MB"},
		{1_204_999_999, "1.1 GB"},
		{1_999_204_999_999, "1.8 TB"},
	}

	for _, test := range tests {
		t.Run(test.out, func(t *testing.T) {
			actualOut := FormatBytes(test.in)
			if actualOut != test.out {
				t.Errorf("FormatBytes(%d) = %q, want %q", test.in, actualOut, test.out)
			}
		})
	}
}
