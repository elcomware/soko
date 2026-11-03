package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

type AuthRequiredMiddleware struct{}

func (m *AuthRequiredMiddleware) Signature() string {
	return "auth_required"
}

func (m *AuthRequiredMiddleware) Handle(ctx http.Context) {
	if !ctx.Request().HasSession() || (ctx.Request().Session().Get("user_id") == nil && ctx.Request().Session().Get("user") == nil) {
		ctx.Request().Session().Flash("error", "Please sign in to access the enterprise dashboard.")
		ctx.Response().Redirect(http.StatusFound, "/login")
		ctx.Request().Abort()
		return
	}
	ctx.Request().Next()
}

func AuthRequired() http.Middleware {
	return &AuthRequiredMiddleware{}
}
