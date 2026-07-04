package ptr

func From[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// To returns a pointer to v. Handy for proto3-optional fields and any API that
// takes *T from a value.
func To[T any](v T) *T {
	return &v
}

// NonZero returns a pointer to v, or nil when v is the zero value. Used to map a
// sentinel-zero scalar (an unset id/count) onto a proto3-optional *T.
func NonZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

func ValueOr[T any](p *T, defaultVal T) T {
	if p == nil {
		return defaultVal
	}
	return *p
}
