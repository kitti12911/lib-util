package ptr

func From[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func ValueOr[T any](p *T, defaultVal T) T {
	if p == nil {
		return defaultVal
	}
	return *p
}
