package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	gobetterauth "github.com/GoBetterAuth/go-better-auth"
	gobetterauthconfig "github.com/GoBetterAuth/go-better-auth/config"
	gobetterauthevents "github.com/GoBetterAuth/go-better-auth/events"
	gobetterauthmodels "github.com/GoBetterAuth/go-better-auth/models"

	"github.com/GoBetterAuth/go-better-auth-playground/events"
	loggerplugin "github.com/GoBetterAuth/go-better-auth-playground/plugins/logger"
	"github.com/GoBetterAuth/go-better-auth-playground/storage"
	"github.com/GoBetterAuth/go-better-auth-playground/utils"
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

	// -------------------------------------
	// Init GoBetterAuth config
	// -------------------------------------

	config := gobetterauthconfig.NewConfig(
		gobetterauthconfig.WithAppName("GoBetterAuthPlayground"),
		gobetterauthconfig.WithBasePath("/api/auth"),
		gobetterauthconfig.WithDatabase(gobetterauthmodels.DatabaseConfig{
			Provider:         "postgres",
			ConnectionString: os.Getenv("DATABASE_URL"),
		}),
		gobetterauthconfig.WithSecondaryStorage(
			gobetterauthmodels.SecondaryStorageConfig{
				Type:    gobetterauthmodels.SecondaryStorageTypeCustom,
				Storage: storage.NewRedisSecondaryStorage(),
			},
		),
		gobetterauthconfig.WithEmailPassword(gobetterauthmodels.EmailPasswordConfig{
			Enabled:                  true,
			DisableSignUp:            false,
			RequireEmailVerification: true,
			AutoSignIn:               true,
			SendResetPasswordEmail: func(user gobetterauthmodels.User, url string, token string) error {
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
		gobetterauthconfig.WithEmailVerification(gobetterauthmodels.EmailVerificationConfig{
			SendOnSignUp: true,
			SendVerificationEmail: func(user gobetterauthmodels.User, url string, token string) error {
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
		gobetterauthconfig.WithUser(gobetterauthmodels.UserConfig{
			ChangeEmail: gobetterauthmodels.ChangeEmailConfig{
				Enabled: true,
				SendEmailChangeVerificationEmail: func(user gobetterauthmodels.User, newEmail string, url string, token string) error {
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
		gobetterauthconfig.WithCSRF(
			gobetterauthmodels.CSRFConfig{
				Enabled: true,
			},
		),
		gobetterauthconfig.WithSocialProviders(
			gobetterauthmodels.SocialProvidersConfig{
				Default: gobetterauthmodels.DefaultOAuth2ProvidersConfig{
					Discord: &gobetterauthmodels.OAuth2Config{
						ClientID:     utils.GetEnv("DISCORD_CLIENT_ID", ""),
						ClientSecret: utils.GetEnv("DISCORD_CLIENT_SECRET", ""),
						RedirectURL:  fmt.Sprintf("%s/api/auth/oauth2/discord/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
					GitHub: &gobetterauthmodels.OAuth2Config{
						ClientID:     utils.GetEnv("GITHUB_CLIENT_ID", ""),
						ClientSecret: utils.GetEnv("GITHUB_CLIENT_SECRET", ""),
						RedirectURL:  fmt.Sprintf("%s/api/auth/oauth2/github/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
					Google: &gobetterauthmodels.OAuth2Config{
						ClientID:     utils.GetEnv("GOOGLE_CLIENT_ID", ""),
						ClientSecret: utils.GetEnv("GOOGLE_CLIENT_SECRET", ""),
						RedirectURL:  fmt.Sprintf("%s/api/auth/oauth2/google/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
				},
			},
		),
		gobetterauthconfig.WithTrustedOrigins(
			gobetterauthmodels.TrustedOriginsConfig{
				Origins: []string{"http://localhost:3000"},
			},
		),
		// Uncomment to test out rate limiting
		// gobetterauthdomain.WithRateLimit(
		// 	gobetterauthdomain.RateLimitConfig{
		// 		Enabled: true,
		// 		Window:  10 * time.Second,
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
		gobetterauthconfig.WithEndpointHooks(
			gobetterauthmodels.EndpointHooksConfig{
				Before: func(ctx *gobetterauthmodels.EndpointHookContext) error {
					logger.Debug(fmt.Sprintf("in 'before' endpoint hook %s %s", ctx.Request.Method, ctx.Request.URL.Path))
					return nil
				},
				// Uncomment this to test out modifying responses
				// Response: func(ctx *gobetterauthmodels.EndpointHookContext) error {
				// 	logger.Debug(fmt.Sprintf("in 'response' endpoint hook %s %s", ctx.Request.Method, ctx.Request.URL.Path))

				// 	if ctx.Path == "/api/protected" {
				// 		ctx.ResponseStatus = http.StatusTeapot
				// 		ctx.ResponseHeaders["Content-Type"] = []string{"text/html"}
				// 		ctx.ResponseBody = []byte("<h1>🍵 I'm a teapot! This response was modified by an endpoint hook.</h1>")
				// 	}

				// 	return nil
				// },
				After: func(ctx *gobetterauthmodels.EndpointHookContext) error {
					logger.Debug(fmt.Sprintf("in 'after' endpoint hook %s %s", ctx.Request.Method, ctx.Request.URL.Path))
					return nil
				},
			},
		),
		gobetterauthconfig.WithDatabaseHooks(gobetterauthmodels.DatabaseHooksConfig{
			Users: &gobetterauthmodels.UserDatabaseHooksConfig{
				BeforeCreate: func(user *gobetterauthmodels.User) error {
					logger.Debug(fmt.Sprintf("in DB hook before creating user with email: %s", user.Email))
					return nil
				},
			},
		}),
		gobetterauthconfig.WithEventHooks(gobetterauthmodels.EventHooksConfig{
			OnUserSignedUp: func(user gobetterauthmodels.User) error {
				logger.Info(fmt.Sprintf("User signed up with email: %s", user.Email))
				return nil
			},
			OnEmailVerified: func(user gobetterauthmodels.User) error {
				logger.Info(fmt.Sprintf("Email verified for user with email: %s", user.Email))
				return nil
			},
			OnEmailChanged: func(user gobetterauthmodels.User) error {
				logger.Info(fmt.Sprintf("User with email %s changed their email", user.Email))
				return nil
			},
			OnPasswordChanged: func(user gobetterauthmodels.User) error {
				logger.Info(fmt.Sprintf("User with email %s changed their password", user.Email))
				return nil
			},
		}),
		gobetterauthconfig.WithEventBus(
			gobetterauthmodels.EventBusConfig{
				Enabled: true,
				Prefix:  "gobetterauthplayground.",
				// Uncomment to test out using watermill with Kafka/Redpanda as the event bus
				PubSub: gobetterauthevents.NewWatermillPubSub(
					events.NewKafkaPublisher(),
					events.NewKafkaSubscriber(),
				),
			},
		),
		gobetterauthconfig.WithPlugins(
			gobetterauthmodels.PluginsConfig{
				Plugins: []gobetterauthmodels.Plugin{
					loggerplugin.NewLoggerPlugin(gobetterauthmodels.PluginConfig{
						Enabled: true,
					}),
				},
			},
		),
	)

	// -------------------------------------
	// Init GoBetterAuth instance and run migrations
	// -------------------------------------

	goBetterAuth := gobetterauth.New(config)
	// You can uncomment the following 2 lines to drop all migrations (i.e., reset the database).
	// goBetterAuth.DropMigrations()
	// return
	goBetterAuth.RunMigrations()

	if id, err := goBetterAuth.EventBus.Subscribe(
		gobetterauthmodels.EventUserSignedUp,
		func(ctx context.Context, event gobetterauthmodels.Event) error {
			var data map[string]any
			if err := json.Unmarshal(event.Payload, &data); err != nil {
				slog.Error("failed to unmarshal json", "error", err)
				return err
			}

			if userId, ok := data["id"]; ok {
				if !ok {
					slog.Error("user ID not found in event payload")
					return fmt.Errorf("user ID not found in event payload")
				}

				slog.Info("EventUserSignedUp event received in main()", "user_id", userId)
			}

			return nil
		}); err != nil {
		slog.Error("failed to subscribe to EventUserSignedUp", "error", err)
		return
	} else {
		slog.Info("subscribed to EventUserSignedUp in main()", "subscription_id", id)
	}

	if id, err := goBetterAuth.EventBus.Subscribe(
		gobetterauthmodels.EventUserLoggedIn,
		func(ctx context.Context, event gobetterauthmodels.Event) error {
			var data map[string]any
			if err := json.Unmarshal(event.Payload, &data); err != nil {
				slog.Error("failed to unmarshal json", "error", err)
				return err
			}

			if userId, ok := data["id"]; ok {
				if !ok {
					slog.Error("user ID not found in event payload")
					return fmt.Errorf("user ID not found in event payload")
				}

				slog.Info("EventUserLoggedIn event received in main()", "user_id", userId)
			}

			return nil
		}); err != nil {
		slog.Error("failed to subscribe to EventUserLoggedIn", "error", err)
		return
	} else {
		slog.Info("subscribed to EventUserLoggedIn in main()", "subscription_id", id)
	}

	// -------------------------------------
	// Choose your api framework of choice to wrap GoBetterAuth (we use Echo in this example)
	// -------------------------------------

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

	// -------------------------------------
	// Custom routes attached to the auth handler for extensibility.
	// This must be done before calling goBetterAuth.Handler()
	// -------------------------------------

	// /api/auth/get-message
	goBetterAuth.RegisterRoute(gobetterauthmodels.CustomRoute{
		Method: "GET",
		Path:   "/get-message",
		Middleware: []gobetterauthmodels.CustomRouteMiddleware{
			goBetterAuth.AuthMiddleware(),
		},
		Handler: func(config *gobetterauthmodels.Config) http.Handler {
			handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(map[string]any{
					"message": fmt.Sprintf("%s: Hello from custom route!", config.AppName),
				})
			})
			return handler
		},
	})

	// /api/auth/send-message
	goBetterAuth.RegisterRoute(gobetterauthmodels.CustomRoute{
		Method: "POST",
		Path:   "/send-message",
		Middleware: []gobetterauthmodels.CustomRouteMiddleware{
			goBetterAuth.AuthMiddleware(),
			goBetterAuth.CSRFMiddleware(),
		},
		Handler: func(config *gobetterauthmodels.Config) http.Handler {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]any{
						"error": "Invalid JSON body",
					})
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(map[string]any{
					"message": fmt.Sprintf("%s: data received", config.AppName),
					"data":    body,
				})
			})
			return handler
		},
	})

	// -------------------------------------

	// -------------------------------------
	// Attach GoBetterAuth handler to your chosen framework and run your server
	// -------------------------------------

	auth.Any("/*", echo.WrapHandler(goBetterAuth.Handler()))

	protected := api.Group(
		"/protected",
		// The order of middleware matters here
		echo.WrapMiddleware(goBetterAuth.CorsAuthMiddleware()),
		echo.WrapMiddleware(goBetterAuth.AuthMiddleware()),
		echo.WrapMiddleware(goBetterAuth.RateLimitMiddleware()),
		echo.WrapMiddleware(goBetterAuth.EndpointHooksMiddleware()),
	)
	protected.GET("", func(c echo.Context) error {
		userId, ok := goBetterAuth.GetUserIDFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"error": "Failed to get user ID from context",
			})
		}
		slog.Debug("User ID from context:", slog.String("user_id", userId))
		return c.JSON(http.StatusOK, map[string]any{
			"userId":  userId,
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
