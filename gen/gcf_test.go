package main

import (
	"testing"
)

type gcftest struct {
	in  []int
	out int
}

var gcftests = []gcftest{
	{[]int{1, 1, 1}, 1},
	{[]int{2, 6}, 2},
	{[]int{20, 140}, 20},
}

func TestGCF(t *testing.T) {
	for _, test := range gcftests {
		if x := gcf(test.in...); x != test.out {
			t.Errorf("got %d expected %d for %v", x, test.out, test.in)
		}
	}
}
