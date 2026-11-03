package controllers

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"

	"elcomware/app/models"
	"elcomware/app/services"
)

type AuthController struct {
	store *services.Store
}

func NewAuthController() *AuthController {
	return &AuthController{
		store: services.GetStore(),
	}
}

func (a *AuthController) ShowLogin(ctx http.Context) http.Response {
	// If already logged in, redirect
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}

	return inertiafacades.Inertia().Render(ctx, "Auth/Login", map[string]any{
		"demoUsers": []map[string]string{
			{
				"role":     "Platform Administrator",
				"email":    "admin@soko.global",
				"password": "admin123",
				"company":  "Soko Global HQ",
				"badge":    "Super Admin",
			},
			{
				"role":     "Enterprise Supplier / Merchant",
				"email":    "merchant@apextech.com",
				"password": "merchant123",
				"company":  "Apex Technologies Group",
				"badge":    "Verified Gold",
			},
			{
				"role":     "Global Procurement Buyer",
				"email":    "buyer@megacorp.com",
				"password": "buyer123",
				"company":  "MegaCorp Procurement Int.",
				"badge":    "Enterprise Buyer",
			},
		},
	})
}

func (a *AuthController) Login(ctx http.Context) http.Response {
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		ctx.Request().Session().Flash("error", "Email and password are required.")
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	user, err := a.store.Authenticate(email, password)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	ctx.Request().Session().Put("user", user)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Welcome back, %s!", user.Name))

	if user.Role == models.RoleAdmin || user.Role == models.RoleMerchant {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}

func (a *AuthController) ShowRegister(ctx http.Context) http.Response {
	if ctx.Request().HasSession() && ctx.Request().Session().Get("user") != nil {
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

	role := models.RoleBuyer
	if accountType == "supplier" {
		role = models.RoleMerchant
	}

	newUser := &models.User{
		Name:        name,
		Email:       email,
		Password:    password,
		Role:        role,
		CompanyName: companyName,
		Country:     country,
		AccountType: accountType,
		Verified:    true,
	}

	created, err := a.store.RegisterUser(newUser)
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusFound, "/register")
	}

	ctx.Request().Session().Put("user", created)
	ctx.Request().Session().Flash("success", fmt.Sprintf("Enterprise account created! Welcome to Soko, %s.", created.Name))

	if created.Role == models.RoleMerchant || created.Role == models.RoleAdmin {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}

func (a *AuthController) Logout(ctx http.Context) http.Response {
	if ctx.Request().HasSession() {
		ctx.Request().Session().Forget("user")
		ctx.Request().Session().Flash("success", "You have successfully signed out.")
	}
	return ctx.Response().Redirect(http.StatusFound, "/")
}

