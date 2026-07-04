package converter

import (
	"math"
	"testing"
)

func TestSlice(t *testing.T) {
	got := Slice([]int{1, 2, 3}, func(n int) int { return n * 2 })
	want := []int{2, 4, 6}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
	// nil/empty input yields a non-nil empty slice.
	if out := Slice(nil, func(n int) int { return n }); out == nil || len(out) != 0 {
		t.Errorf("Slice(nil) = %v, want non-nil empty", out)
	}
}

func TestSlicePtr(t *testing.T) {
	type row struct{ n int }
	rows := []row{{1}, {2}}
	got := SlicePtr(rows, func(r *row) int { return r.n * 10 })
	if len(got) != 2 || got[0] != 10 || got[1] != 20 {
		t.Fatalf("SlicePtr = %v, want [10 20]", got)
	}
}

func TestSliceDeref(t *testing.T) {
	// Odd numbers map to nil and are skipped; even numbers are dereferenced.
	got := SliceDeref([]int{1, 2, 3, 4}, func(n int) *int {
		if n%2 != 0 {
			return nil
		}
		v := n
		return &v
	})
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatalf("SliceDeref = %v, want [2 4]", got)
	}
}

func TestInt32FromInt(t *testing.T) {
	cases := map[int]int32{
		5:                 42 - 37,
		math.MaxInt32:     math.MaxInt32,
		math.MaxInt32 + 1: math.MaxInt32,
		math.MinInt32 - 1: math.MinInt32,
	}
	for in, want := range cases {
		if got := Int32FromInt(in); got != want {
			t.Errorf("Int32FromInt(%d) = %d, want %d", in, got, want)
		}
	}
}
