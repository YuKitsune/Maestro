package handlers

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type Result[T model.Thing] struct {
	Type  model.Type
	Items []T
}

func NewResult[T model.Thing](typ model.Type) *Result[T] {
	return &Result[T]{
		Type:  typ,
		Items: []T{},
	}
}

func (r *Result[T]) Add(t T) {
	if r.HasResultFor(t.GetSource()) {
		for i, item := range r.Items {
			if item.GetSource() == t.GetSource() {
				r.Items[i] = t
			}
		}
	} else {
		r.Items = append(r.Items, t)
	}
}

func (r *Result[T]) AddAll(ts []T) {
	for _, t := range ts {
		r.Add(t)
	}
}

func (r *Result[T]) HasResultFor(key model.StreamingServiceType) bool {
	for _, v := range r.Items {
		if v.GetSource() == key {
			return true
		}
	}

	return false
}

func (r *Result[T]) HasResults() bool {
	return len(r.Items) > 0
}
