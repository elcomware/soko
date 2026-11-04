package routes

import (
	"elcomware/app/http/controllers"
	appmiddleware "elcomware/app/http/middleware"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	sessionmiddleware "github.com/goravel/framework/session/middleware"
)

func Web() {
	facades.Route().Static("public", "./public")
	facades.Route().Static("build", "./public/build")
	facades.Route().Middleware(sessionmiddleware.StartSession(), appmiddleware.SameOrigin(), appmiddleware.HandleInertiaRequests()).Group(func(router route.Router) {
		store := controllers.NewStoreController()
		auth := controllers.NewAuthController()
		dashboard := controllers.NewDashboardController()
		commerce := controllers.NewCommerceController()
		router.Get("/", store.Index)
		router.Get("/products/{id}", store.Show)
		router.Get("/cart", commerce.Cart)
		router.Post("/cart/items", commerce.AddItem)
		router.Post("/cart/items/{id}", commerce.UpdateItem)
		router.Post("/cart/items/{id}/delete", commerce.RemoveItem)
		router.Get("/login", auth.ShowLogin)
		router.Post("/login", auth.Login)
		router.Get("/register", auth.ShowRegister)
		router.Post("/register", auth.Register)
		router.Post("/logout", auth.Logout)
		router.Middleware(appmiddleware.AuthRequired()).Group(func(r route.Router) {
			r.Post("/checkout", commerce.Checkout)
			r.Get("/checkout", commerce.Cart)
			r.Post("/make-offer", store.MakeOffer)
			r.Post("/rfqs", store.CreateRfq)
			r.Get("/account", commerce.Account)
			r.Post("/account/orders/{id}/cancel", commerce.CancelOrder)
			r.Post("/account/rfqs/{id}/quotes/{quoteId}/accept", commerce.AcceptQuote)
			r.Post("/account/rfqs/{id}/close", commerce.CloseRfq)
			r.Get("/dashboard", dashboard.Index)
			r.Post("/dashboard/products", dashboard.CreateProduct)
			r.Post("/dashboard/products/{id}", dashboard.UpdateProduct)
			r.Post("/dashboard/products/{id}/delete", dashboard.DeleteProduct)
			r.Post("/dashboard/rfqs/quote", dashboard.SubmitRfqQuote)
			r.Post("/dashboard/offers/{id}/decision", dashboard.DecideOffer)
			r.Post("/dashboard/content", controllers.NewContentController().Save)
		})
	})
}
