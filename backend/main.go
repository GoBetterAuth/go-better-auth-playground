package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	gobetterauth "github.com/GoBetterAuth/go-better-auth"
	"github.com/GoBetterAuth/go-better-auth-playground/storage"
	"github.com/GoBetterAuth/go-better-auth-playground/utils"
	gobetterauthdomain "github.com/GoBetterAuth/go-better-auth/pkg/domain"
)

// Feel free to change this implementation to use your own mailer service e.g. SendGrid/Resend etc.
func sendEmail(to string, subject string, html string, text string) error {
	payload := map[string]any{
		"from":    utils.GetEnv("MAILER_FROM_ADDRESS", ""),
		"to":      []string{to},
		"subject": subject,
		"html":    html,
		"text":    text,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(utils.GetEnv("MAILER_URL", ""), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Init GoBetterAuth config

	config := gobetterauthdomain.NewConfig(
		gobetterauthdomain.WithAppName("GoBetterAuthPlayground"),
		gobetterauthdomain.WithBasePath("/api/auth"),
		gobetterauthdomain.WithDatabase(gobetterauthdomain.DatabaseConfig{
			Provider:         "postgres",
			ConnectionString: os.Getenv("DATABASE_URL"),
		}),
		gobetterauthdomain.WithSecondaryStorage(
			gobetterauthdomain.SecondaryStorageConfig{
				Type:    gobetterauthdomain.SecondaryStorageTypeCustom,
				Storage: storage.NewRedisSecondaryStorage(),
			},
		),
		gobetterauthdomain.WithEmailPassword(gobetterauthdomain.EmailPasswordConfig{
			Enabled:                  true,
			DisableSignUp:            false,
			RequireEmailVerification: true,
			AutoSignIn:               true,
			SendResetPasswordEmail: func(user gobetterauthdomain.User, url, token string) error {
				if err := sendEmail(
					user.Email,
					"Reset your password",
					fmt.Sprintf("<p>Please reset your password by clicking <a href=\"%s\">here</a>.</p>", url),
					fmt.Sprintf("Please reset your password by visiting the following link: %s", url),
				); err != nil {
					return err
				}
				return nil
			},
		}),
		gobetterauthdomain.WithEmailVerification(gobetterauthdomain.EmailVerificationConfig{
			SendOnSignUp: true,
			SendVerificationEmail: func(user gobetterauthdomain.User, url string, token string) error {
				if err := sendEmail(
					user.Email,
					"Verify your email",
					fmt.Sprintf("<p>Please verify your email by clicking <a href=\"%s\">here</a>.</p>", url),
					fmt.Sprintf("Please verify your email by visiting the following link: %s", url),
				); err != nil {
					return err
				}
				return nil
			},
		}),
		gobetterauthdomain.WithUser(gobetterauthdomain.UserConfig{
			ChangeEmail: gobetterauthdomain.ChangeEmailConfig{
				Enabled: true,
				SendEmailChangeVerificationEmail: func(user gobetterauthdomain.User, newEmail string, url string, token string) error {
					if err := sendEmail(
						user.Email,
						"You requested to change your email",
						fmt.Sprintf("<p>Please click on the following link to change your email from %s to %s <a href=\"%s\">here</a>.</p>", user.Email, newEmail, url),
						fmt.Sprintf("Please click on the following link to change your email from %s to %s: %s", user.Email, newEmail, url),
					); err != nil {
						return err
					}
					return nil
				},
			},
		}),
		gobetterauthdomain.WithCSRF(
			gobetterauthdomain.CSRFConfig{
				Enabled: true,
			},
		),
		gobetterauthdomain.WithSocialProviders(
			gobetterauthdomain.SocialProvidersConfig{
				Default: gobetterauthdomain.DefaultOAuth2ProvidersConfig{
					Discord: &gobetterauthdomain.OAuth2Config{
						ClientID:     utils.GetEnv("DISCORD_CLIENT_ID", ""),
						ClientSecret: utils.GetEnv("DISCORD_CLIENT_SECRET", ""),
						RedirectURL:  fmt.Sprintf("%s/api/auth/oauth2/discord/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
					GitHub: &gobetterauthdomain.OAuth2Config{
						ClientID:     utils.GetEnv("GITHUB_CLIENT_ID", ""),
						ClientSecret: utils.GetEnv("GITHUB_CLIENT_SECRET", ""),
						RedirectURL:  fmt.Sprintf("%s/api/auth/oauth2/github/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
					Google: &gobetterauthdomain.OAuth2Config{
						ClientID:     utils.GetEnv("GOOGLE_CLIENT_ID", ""),
						ClientSecret: utils.GetEnv("GOOGLE_CLIENT_SECRET", ""),
						RedirectURL:  fmt.Sprintf("%s/api/auth/oauth2/google/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
				},
			},
		),
		gobetterauthdomain.WithTrustedOrigins(
			gobetterauthdomain.TrustedOriginsConfig{
				Origins: []string{"http://localhost:3000"},
			},
		),
		// Uncomment to test out rate limiting
		// gobetterauthdomain.WithRateLimit(
		// 	gobetterauthdomain.RateLimitConfig{
		// 		Enabled: true,
		// 		Window:  30 * time.Second,
		// 		Max:     5,
		// 		CustomRules: map[string]gobetterauthdomain.RateLimitCustomRuleFunc{
		// 			"/api/protected": func(req *http.Request) gobetterauthdomain.RateLimitCustomRule {
		// 				return gobetterauthdomain.RateLimitCustomRule{
		// 					Disabled: true,
		// 				}
		// 			},
		// 		},
		// 	},
		// ),
		gobetterauthdomain.WithEndpointHooks(
			gobetterauthdomain.EndpointHooksConfig{
				Before: func(ctx *gobetterauthdomain.EndpointHookContext) error {
					logger.Debug(fmt.Sprintf("in endpoint hook before %s %s", ctx.Request.Method, ctx.Request.URL.Path))
					return nil
				},
				After: func(ctx *gobetterauthdomain.EndpointHookContext) error {
					logger.Debug(fmt.Sprintf("in endpoint hook after %s %s", ctx.Request.Method, ctx.Request.URL.Path))
					return nil
				},
			},
		),
		gobetterauthdomain.WithDatabaseHooks(gobetterauthdomain.DatabaseHooksConfig{
			Users: &gobetterauthdomain.UserDatabaseHooksConfig{
				BeforeCreate: func(user *gobetterauthdomain.User) error {
					logger.Debug(fmt.Sprintf("in DB hook before creating user with email: %s", user.Email))
					return nil
				},
			},
		}),
		gobetterauthdomain.WithEventHooks(gobetterauthdomain.EventHooksConfig{
			OnUserSignedUp: func(user gobetterauthdomain.User) error {
				logger.Info(fmt.Sprintf("User signed up with email: %s", user.Email))
				return nil
			},
			OnEmailVerified: func(user gobetterauthdomain.User) error {
				logger.Info(fmt.Sprintf("Email verified for user with email: %s", user.Email))
				return nil
			},
			OnEmailChanged: func(user gobetterauthdomain.User) error {
				logger.Info(fmt.Sprintf("User with email %s changed their email", user.Email))
				return nil
			},
			OnPasswordChanged: func(user gobetterauthdomain.User) error {
				logger.Info(fmt.Sprintf("User with email %s changed their password", user.Email))
				return nil
			},
		}),
	)

	// Init GoBetterAuth instance and run migrations

	goBetterAuth := gobetterauth.New(config, nil)
	// You can uncomment the following 2 lines to drop all migrations (i.e., reset the database).
	// goBetterAuth.DropMigrations()
	// return
	goBetterAuth.RunMigrations()

	// Choose your api framework of choice to wrap GoBetterAuth (we use Echo in this example)

	echoInstance := echo.New()
	if err != nil {
		echoInstance.Logger.Fatal(err)
	}

	api := echoInstance.Group("/api")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Hello, World!",
		})
	})

	auth := api.Group(
		"/auth",
		echo.WrapMiddleware(goBetterAuth.CorsAuthMiddleware()),
		echo.WrapMiddleware(goBetterAuth.OptionalAuthMiddleware()),
	)
	auth.Any("/*", echo.WrapHandler(goBetterAuth.Handler()))

	protected := api.Group(
		"/protected",
		echo.WrapMiddleware(goBetterAuth.CorsAuthMiddleware()),
		echo.WrapMiddleware(goBetterAuth.AuthMiddleware()),
		echo.WrapMiddleware(goBetterAuth.RateLimitMiddleware()),
	)
	protected.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Protected Route!",
		})
	})
	// This route is protected by CSRF so requires csrf token cookie and header
	protected.POST("/data", func(c echo.Context) error {
		var body map[string]any
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": "Invalid JSON body",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": body,
		})
	}, echo.WrapMiddleware(goBetterAuth.CSRFMiddleware()))

	echoInstance.Logger.Fatal(echoInstance.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
