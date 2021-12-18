package middleware

import (
	"context"
	"github.com/yukitsune/camogo"
	maestroContext "maestro/pkg/api/context"
	"maestro/pkg/api/handlers"
	"net/http"
)

type ContextInjection struct {
	container camogo.Container
}

func NewContainerInjectionMiddleware(container camogo.Container) *ContextInjection {
	return &ContextInjection{container: container}
}

func (m *ContextInjection) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Behaviour here is kinda weird:
		// - ctx needs the child container
		// - child container needs a context

		// So:
		// - The container gets the original context (without the container itself inside)
		// - The request gets the new context with the container

		// Meaning:
		// - No circular reference (ctr -> ctx -> ctr -> ctx -> ...)
		// - Container access from context limited to HTTP handlers

		cc, err := m.container.NewChildWith(func(cb camogo.ContainerBuilder) error {
			err := cb.RegisterFactory(func() context.Context {
				return r.Context()
			}, camogo.SingletonLifetime)

			return err
		})

		if err != nil {
			handlers.Error(w, err)
			return
		}

		ctx := maestroContext.WithContainer(r.Context(), cc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
