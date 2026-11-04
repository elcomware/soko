package controllers

import (
	"fmt"
	frameworksession "github.com/goravel/framework/session"
	"strings"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	appmiddleware "elcomware/app/http/middleware"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (a *AuthController) ShowLogin(ctx http.Context) http.Response {
	if appmiddleware.CurrentUser(ctx) != nil {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	return inertiafacades.Inertia().Render(ctx, "Auth/Login", map[string]any{
		"demoUsers": []map[string]string{
			{
				"role":     "Platform Administrator",
				"email":    "admin@soko.global",
				"password": "secret123",
				"company":  "Soko Multi-Region Marketplace Operations Inc.",
				"badge":    "Super Admin",
			},
			{
				"role":     "Enterprise Supplier / Merchant",
				"email":    "supplier.ea@soko.global",
				"password": "secret123",
				"company":  "Apex Heavy Agro-Tech & Grain Silos Ltd",
				"badge":    "Verified Gold OEM",
			},
			{
				"role":     "Global Procurement Buyer",
				"email":    "buyer@soko.global",
				"password": "secret123",
				"company":  "Great Rift Infrastructure Consortium",
				"badge":    "Enterprise Buyer",
			},
		},
	})
}

func (a *AuthController) Login(ctx http.Context) http.Response {
	var req requests.LoginRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil || errors != nil {
		errorMsg := "Valid email and password are required."
		if errors != nil {
			errorMsg = errors.One()
		}
		ctx.Request().Session().Flash("error", errorMsg)
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	user, err := a.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	if err := ctx.Request().Session().Regenerate(true); err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "Could not start a secure session.")
	}
	frameworksession.WriteCookie(ctx, ctx.Request().Session())
	ctx.Request().Session().Forget("user")
	ctx.Request().Session().Put("user_id", user.ID)
	ctx.Request().Session().Put("user_email", user.Email)
	ctx.Request().Session().Put("user_name", user.Name)
	ctx.Request().Session().Put("user_role", string(user.Role))
	ctx.Request().Session().Flash("success", fmt.Sprintf("Welcome back, %s!", user.Name))

	if user.Role == models.RoleAdmin || user.Role == models.RoleMerchant {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}
	return ctx.Response().Redirect(http.StatusFound, "/account")
}

func (a *AuthController) ShowRegister(ctx http.Context) http.Response {
	if appmiddleware.CurrentUser(ctx) != nil {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	return inertiafacades.Inertia().Render(ctx, "Auth/Register", map[string]any{})
}

func (a *AuthController) Register(ctx http.Context) http.Response {
	var req requests.RegisterRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil || errors != nil {
		errorMsg := "Please check your registration details."
		if errors != nil {
			errorMsg = errors.One()
		}
		ctx.Request().Session().Flash("error", errorMsg)
		return ctx.Response().Redirect(http.StatusFound, "/register")
	}

	role := models.RoleBuyer
	if req.AccountType == "supplier" {
		role = models.RoleMerchant
	}

	accountType := strings.TrimSpace(req.AccountType)
	if accountType == "" {
		if role == models.RoleMerchant {
			accountType = "supplier"
		} else {
			accountType = "individual"
		}
	}

	user := &models.User{
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		Role:        role,
		CompanyName: req.CompanyName,
		Country:     req.Country,
		AccountType: accountType,
	}

	created, err := a.authService.Register(user)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/register")
	}

	if err := ctx.Request().Session().Regenerate(true); err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "Could not start a secure session.")
	}
	frameworksession.WriteCookie(ctx, ctx.Request().Session())
	ctx.Request().Session().Forget("user")
	ctx.Request().Session().Put("user_id", created.ID)
	ctx.Request().Session().Put("user_email", created.Email)
	ctx.Request().Session().Put("user_name", created.Name)
	ctx.Request().Session().Put("user_role", string(created.Role))
	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise account created! Welcome to Soko, %s.", created.Name))

	if created.Role == models.RoleMerchant || created.Role == models.RoleAdmin {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}
	return ctx.Response().Redirect(http.StatusFound, "/account")
}

func (a *AuthController) Logout(ctx http.Context) http.Response {
	if ctx.Request().HasSession() {
		if err := ctx.Request().Session().Invalidate(); err != nil {
			return ctx.Response().String(http.StatusInternalServerError, "Could not end the session. Please try again.")
		}
		frameworksession.WriteCookie(ctx, ctx.Request().Session())
		ctx.Request().Session().Forget("user", "user_id")
		ctx.Request().Session().Forget("user_email")
		ctx.Request().Session().Forget("user_name")
		ctx.Request().Session().Forget("user_role")
		ctx.Request().Session().Flash("success", "You have successfully signed out.")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}
