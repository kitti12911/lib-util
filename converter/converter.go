// Package converter holds small, dependency-free helpers shared by the service
// mapping layers: generic slice mappers that collapse the repetitive
// "make + range + append" loops, plus a few primitive conversions that were
// otherwise copy-pasted per feature package.
package converter

import "math"

// Slice maps each element of src through fn, preserving order and length. A nil
// or empty src yields a non-nil empty slice, matching the hand-written
// `make([]D, len(src))` loops it replaces.
func Slice[S, D any](src []S, fn func(S) D) []D {
	out := make([]D, len(src))
	for i, s := range src {
		out[i] = fn(s)
	}
	return out
}

// SlicePtr is like Slice but passes a pointer to each element, for mappers with
// a `func(*T) D` signature (e.g. proto row mappers that take `*Row`). It indexes
// src so the pointer refers to the backing array element, not a loop copy.
func SlicePtr[S, D any](src []S, fn func(*S) D) []D {
	out := make([]D, len(src))
	for i := range src {
		out[i] = fn(&src[i])
	}
	return out
}

// SliceDeref maps each element through fn, which returns *D, skipping nil results
// and dereferencing the rest. It suits response wrappers that turn a proto
// repeated field into a value slice via a pointer-returning element mapper.
func SliceDeref[S, D any](src []S, fn func(S) *D) []D {
	out := make([]D, 0, len(src))
	for _, s := range src {
		if d := fn(s); d != nil {
			out = append(out, *d)
		}
	}
	return out
}

// Int32FromInt narrows an int to int32, saturating at the int32 bounds instead
// of wrapping. Replaces the per-package int32FromInt copies.
func Int32FromInt(v int) int32 {
	switch {
	case v > math.MaxInt32:
		return math.MaxInt32
	case v < math.MinInt32:
		return math.MinInt32
	default:
		return int32(v)
	}
}
