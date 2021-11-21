package context

import (
	"context"
	"fmt"
	"github.com/yukitsune/camogo"
)

const ContainerKey = "maestro.container"

func WithContainer(ctx context.Context, container camogo.Container) context.Context {
	return context.WithValue(ctx, ContainerKey, container)
}

func Container(ctx context.Context) (camogo.Container, error) {
	val := ctx.Value(ContainerKey)
	if val != nil {
		container, ok := val.(camogo.Container)
		if ok {
			return container, nil
		}
	}

	return nil, fmt.Errorf("could not find container with key \"%s\" in context", ContainerKey)
}
