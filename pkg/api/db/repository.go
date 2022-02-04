package db

import (
	"context"
	"github.com/yukitsune/maestro/pkg/model"
)

type Repository interface {
	GetThingByLink(context.Context, string) (model.Thing, error)
	GetThingsByGroupId(context.Context, model.ThingGroupId) ([]model.Thing, error)
	AddThings(context.Context, []model.Thing) (int, error)
}
