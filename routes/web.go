package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	sessionmiddleware "github.com/goravel/framework/session/middleware"

	"elcomware/app/http/controllers"
	appmiddleware "elcomware/app/http/middleware"
)

func Web() {
	// Static assets (no session needed).
	facades.Route().Static("public", "./public")
	facades.Route().Static("build", "./public/build")

	// Inertia web group with session support
	facades.Route().Middleware(
		sessionmiddleware.StartSession(),
		appmiddleware.HandleInertiaRequests(),
	).Group(func(router route.Router) {
		storeController := controllers.NewStoreController()
		authController := controllers.NewAuthController()
		dashboardController := controllers.NewDashboardController()

		// eStore Landing, Global Catalog & B2B Features
		router.Get("/", storeController.Index)
		router.Get("/products/{id}", storeController.Show)
		router.Post("/checkout", storeController.Checkout)
		router.Post("/make-offer", storeController.MakeOffer)
		router.Post("/rfqs", storeController.CreateRfq)

		// Authentication & Account Creation
		router.Get("/login", authController.ShowLogin)
		router.Post("/login", authController.Login)
		router.Get("/register", authController.ShowRegister)
		router.Post("/register", authController.Register)
		router.Any("/logout", authController.Logout)

		// Enterprise Admin & Merchant Dashboard
		router.Middleware(appmiddleware.AuthRequired()).Group(func(authRouter route.Router) {
			authRouter.Get("/dashboard", dashboardController.Index)
			authRouter.Post("/dashboard/products", dashboardController.CreateProduct)
			authRouter.Post("/dashboard/products/{id}/delete", dashboardController.DeleteProduct)
			authRouter.Post("/dashboard/rfqs/quote", dashboardController.SubmitRfqQuote)
		})
	})
}
