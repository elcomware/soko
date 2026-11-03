package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"

	"elcomware/config"
	"elcomware/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithSeeders(Seeders).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithMiddleware(func(handler configuration.Middleware) {
			// Global middleware stack (Goravel v1.18)
			// Route-specific session & inertia middleware are scoped within routes.Web()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
