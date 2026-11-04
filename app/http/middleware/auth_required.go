package middleware

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
	"github.com/goravel/framework/contracts/http"
)

func CurrentUser(ctx http.Context) *models.User {
	if !ctx.Request().HasSession() {
		return nil
	}
	id, ok := ctx.Request().Session().Get("user_id").(string)
	if !ok || id == "" {
		return nil
	}
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", id).First(&user); err != nil || user.ID == "" {
		return nil
	}
	queries.HydrateAccess(&user)
	return &user
}

type AuthRequiredMiddleware struct{}

func (m *AuthRequiredMiddleware) Signature() string { return "auth_required" }
func (m *AuthRequiredMiddleware) Handle(ctx http.Context) {
	if CurrentUser(ctx) == nil {
		if ctx.Request().HasSession() {
			ctx.Request().Session().Forget("user_id", "user", "user_name", "user_email", "user_role")
			ctx.Request().Session().Flash("error", "Please sign in to continue.")
		}
		_ = ctx.Response().Redirect(http.StatusFound, "/login").Abort()
		return
	}
	ctx.Request().Next()
}
func AuthRequired() http.Middleware { return &AuthRequiredMiddleware{} }
