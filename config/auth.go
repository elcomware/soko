package config

import (
	"elcomware/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("auth", map[string]any{
		// Authentication Defaults
		//
		// This option controls the default authentication "guard" for your application.
		// "web" uses the session driver for web routes/Inertia, while "api" uses JWT.
		"defaults": map[string]any{
			"guard": "web",
		},

		// Authentication Guards
		//
		// Supported drivers: "jwt", "session"
		"guards": map[string]any{
			"web": map[string]any{
				"driver":   "session",
				"provider": "user",
			},
			"api": map[string]any{
				"driver":   "jwt",
				"provider": "user",
			},
		},

		// User Providers
		// Supported: "orm"
		"providers": map[string]any{
			"user": map[string]any{
				"driver": "orm",
			},
		},
	})
}
