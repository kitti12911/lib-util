package huma

import (
	"bytes"
	"encoding/json"
	"reflect"

	basehuma "github.com/danielgtaylor/huma/v2"
)

type Patch[T any] struct {
	Set   bool
	Null  bool
	Value T
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	p.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		p.Null = true
		var zero T
		p.Value = zero
		return nil
	}
	p.Null = false
	return json.Unmarshal(data, &p.Value)
}

func (p Patch[T]) MarshalJSON() ([]byte, error) {
	if !p.Set || p.Null {
		return []byte("null"), nil
	}
	return json.Marshal(p.Value)
}

func (p Patch[T]) Schema(r basehuma.Registry) *basehuma.Schema {
	var zero T
	t := reflect.TypeOf(zero)
	if t == nil {
		return &basehuma.Schema{Nullable: true}
	}
	s := basehuma.SchemaFromType(r, t)
	s.Nullable = true
	return s
}
