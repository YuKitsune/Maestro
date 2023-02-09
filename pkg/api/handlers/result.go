package handlers

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type Result struct {
	Type  model.Type
	Items []model.HasStreamingService
}

func NewResult(typ model.Type) *Result {
	return &Result{typ, []model.HasStreamingService{}}
}

func (r *Result) Add(v model.HasStreamingService) {
	if r.HasResultFor(v.GetSource()) {
		for i, item := range r.Items {
			if item.GetSource() == v.GetSource() {
				r.Items[i] = v
			}
		}
	} else {
		r.Items = append(r.Items, v)
	}
}

func (r *Result) AddAll(vs []model.HasStreamingService) {
	for _, v := range vs {
		r.Add(v)
	}
}

func (r *Result) HasResultFor(key model.StreamingServiceType) bool {
	for _, v := range r.Items {
		if v.GetSource() == key {
			return true
		}
	}

	return false
}

func (r *Result) HasResults() bool {
	return len(r.Items) > 0
}
