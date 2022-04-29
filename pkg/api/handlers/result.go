package handlers

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type Result struct {
	Type     model.Type
	Services map[model.StreamingServiceKey]interface{}
}

func NewResult(typ model.Type) *Result {
	svcMap := make(map[model.StreamingServiceKey]interface{})
	return &Result{typ, svcMap}
}

func (r *Result) Add(v model.HasStreamingService) {
	r.Services[v.GetSource()] = v
}

func (r *Result) AddAll(vs []model.HasStreamingService) {
	for _, v := range vs {
		r.Add(v)
	}
}

func (r *Result) IsMissingResults() bool {
	for _, v := range r.Services {
		if v == nil {
			return true
		}
	}

	return false
}

func (r *Result) HasResultFor(key model.StreamingServiceKey) bool {
	for k := range r.Services {
		if k == key {
			return true
		}
	}

	return false
}

func (r *Result) HasResults() bool {
	for _, v := range r.Services {
		if v != nil {
			return true
		}
	}

	return false
}
