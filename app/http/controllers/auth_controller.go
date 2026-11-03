package controllers

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	"elcomware/app/dto"
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
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user_id") != nil {
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
	email := strings.TrimSpace(ctx.Request().Input("email"))
	password := strings.TrimSpace(ctx.Request().Input("password"))

	if email == "" || password == "" {
		ctx.Request().Session().Flash("error", "Email and password are required.")
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	input := dto.LoginInput{
		Email:    email,
		Password: password,
	}

	user, err := a.authService.Login(input)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	ctx.Request().Session().Put("user_id", user.ID)
	ctx.Request().Session().Put("user_email", user.Email)
	ctx.Request().Session().Put("user_name", user.Name)
	ctx.Request().Session().Put("user_role", user.Role)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Welcome back, %s!", user.Name))

	if user.Role == "admin" || user.Role == "merchant" {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}

func (a *AuthController) ShowRegister(ctx http.Context) http.Response {
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user_id") != nil {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	return inertiafacades.Inertia().Render(ctx, "Auth/Register", map[string]any{})
}

func (a *AuthController) Register(ctx http.Context) http.Response {
	name := strings.TrimSpace(ctx.Request().Input("name"))
	email := strings.TrimSpace(ctx.Request().Input("email"))
	password := strings.TrimSpace(ctx.Request().Input("password"))
	companyName := strings.TrimSpace(ctx.Request().Input("company_name"))
	country := strings.TrimSpace(ctx.Request().Input("country"))
	accountType := strings.TrimSpace(ctx.Request().Input("account_type")) // "supplier" or "buyer"

	if name == "" || email == "" || password == "" {
		ctx.Request().Session().Flash("error", "Name, email, and password are required.")
		return ctx.Response().Redirect(http.StatusFound, "/register")
	}

	input := dto.RegisterInput{
		Name:        name,
		Email:       email,
		Password:    password,
		CompanyName: companyName,
		Country:     country,
		AccountType: accountType,
	}

	created, err := a.authService.Register(input)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/register")
	}

	ctx.Request().Session().Put("user_id", created.ID)
	ctx.Request().Session().Put("user_email", created.Email)
	ctx.Request().Session().Put("user_name", created.Name)
	ctx.Request().Session().Put("user_role", created.Role)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise account created! Welcome to Soko, %s.", created.Name))

	if created.Role == "merchant" || created.Role == "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}

func (a *AuthController) Logout(ctx http.Context) http.Response {
	if ctx.Request().HasSession() {
		ctx.Request().Session().Forget("user_id")
		ctx.Request().Session().Forget("user_email")
		ctx.Request().Session().Forget("user_name")
		ctx.Request().Session().Forget("user_role")
		ctx.Request().Session().Flash("success", "You have successfully signed out.")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}
