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
	gobetterauthmodels "github.com/GoBetterAuth/go-better-auth/models"

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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
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
				go func() {
					if err := sendEmail(
						user.Email,
						"Reset your password",
						fmt.Sprintf("<p>Please reset your password by clicking <a href=\"%s\">here</a>.</p>", url),
						fmt.Sprintf("Please reset your password by visiting the following link: %s", url),
					); err != nil {
						fmt.Println(err.Error())
					}
				}()

				return nil
			},
		}),
		gobetterauthconfig.WithEmailVerification(gobetterauthmodels.EmailVerificationConfig{
			SendOnSignUp: true,
			SendVerificationEmail: func(user gobetterauthmodels.User, url string, token string) error {
				go func() {
					if err := sendEmail(
						user.Email,
						"Verify your email",
						fmt.Sprintf("<p>Please verify your email by clicking <a href=\"%s\">here</a>.</p>", url),
						fmt.Sprintf("Please verify your email by visiting the following link: %s", url),
					); err != nil {
						fmt.Println(err.Error())
					}
				}()
				return nil
			},
		}),
		gobetterauthconfig.WithUser(gobetterauthmodels.UserConfig{
			ChangeEmail: gobetterauthmodels.ChangeEmailConfig{
				Enabled: true,
				SendEmailChangeVerificationEmail: func(user gobetterauthmodels.User, newEmail string, url string, token string) error {
					go func() {
						if err := sendEmail(
							user.Email,
							"You requested to change your email",
							fmt.Sprintf("<p>Please click on the following link to change your email from %s to %s <a href=\"%s\">here</a>.</p>", user.Email, newEmail, url),
							fmt.Sprintf("Please click on the following link to change your email from %s to %s: %s", user.Email, newEmail, url),
						); err != nil {
							fmt.Println(err.Error())
						}
					}()

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
				Providers: map[string]gobetterauthmodels.OAuth2ProviderConfig{
					"discord": {
						Enabled:     true,
						RedirectURL: fmt.Sprintf("%s/api/auth/oauth2/discord/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
					"github": {
						Enabled:     true,
						RedirectURL: fmt.Sprintf("%s/api/auth/oauth2/github/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
					},
					"google": {
						Enabled:     true,
						RedirectURL: fmt.Sprintf("%s/api/auth/oauth2/google/callback", utils.GetEnv("GO_BETTER_AUTH_BASE_URL", "")),
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
		// gobetterauthconfig.WithRateLimit(
		// 	gobetterauthmodels.RateLimitConfig{
		// 		Enabled: true,
		// 		Window:  30 * time.Second,
		// 		Max:     5,
		// 		CustomRules: map[string]gobetterauthmodels.RateLimitCustomRule{
		// 			"/api/protected": {
		// 				Disabled: true,
		// 			},
		// 		},
		// 	},
		// ),
		gobetterauthconfig.WithEndpointHooks(
			gobetterauthmodels.EndpointHooksConfig{
				Before: func(ctx *gobetterauthmodels.EndpointHookContext) error {
					logger.Info(fmt.Sprintf("in 'before' endpoint hook %s %s", ctx.Method, ctx.Path))

					// Uncomment this to test out custom validation before a handler is executed
					// if ctx.Path == "/api/auth/sign-up/email" {
					// 	if ctx.Method != "POST" {
					// 		return fmt.Errorf("only POST is allowed for %s", ctx.Path)
					// 	}

					// 	email, ok := ctx.Body["email"].(string)
					// 	if !ok {
					// 		return fmt.Errorf("email is required in request body")
					// 	}

					// 	if !strings.HasSuffix(email, "@gmail.com") {
					// 		// Short-circuit with custom error response
					// 		ctx.ResponseHeaders["Content-Type"] = []string{"application/json"}
					// 		ctx.ResponseStatus = http.StatusForbidden
					// 		ctx.ResponseBody = []byte(`{"message": "Only @gmail.com email addresses are allowed to sign up."}`)
					// 		ctx.Handled = true
					// 		// Do not return an error here
					// 		return nil
					// 	}
					// }

					return nil
				},
				// Uncomment this to test out modifying responses
				// Response: func(ctx *gobetterauthmodels.EndpointHookContext) error {
				// 	logger.Debug(fmt.Sprintf("in 'response' endpoint hook %s %s", ctx.Method, ctx.Path))

				// 	if ctx.Path == "/api/protected" {
				// 		ctx.ResponseStatus = http.StatusTeapot
				// 		ctx.ResponseHeaders["Content-Type"] = []string{"text/html"}
				// 		ctx.ResponseBody = []byte("<h1>🍵 I'm a teapot! This response was modified by an endpoint hook.</h1>")
				// 	}

				// 	return nil
				// },
				After: func(ctx *gobetterauthmodels.EndpointHookContext) {
					logger.Debug(fmt.Sprintf("in 'after' endpoint hook %s %s", ctx.Method, ctx.Path))
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
			OnUserSignedUp: func(user gobetterauthmodels.User) {
				logger.Info(fmt.Sprintf("User signed up with email: %s", user.Email))
			},
			OnEmailVerified: func(user gobetterauthmodels.User) {
				logger.Info(fmt.Sprintf("Email verified for user with email: %s", user.Email))
			},
			OnEmailChanged: func(user gobetterauthmodels.User) {
				logger.Info(fmt.Sprintf("User with email %s changed their email", user.Email))
			},
			OnPasswordChanged: func(user gobetterauthmodels.User) {
				logger.Info(fmt.Sprintf("User with email %s changed their password", user.Email))
			},
		}),
		// You can uncomment the following to test out webhooks.
		// Would recommend using event hooks instead of webhooks for these types of hooks
		// as it is more efficient since you can run the logic directly in Go code instead of sending a HTTP request.
		// gobetterauthconfig.WithWebhooks(
		// 	gobetterauthmodels.WebhooksConfig{
		// 		OnUserSignedUp: &gobetterauthmodels.WebhookConfig{
		// 			URL: "<webhook url>",
		// 			Headers: map[string]string{
		// 				"X-API-KEY": "your-api-key",
		// 			},
		// 		},
		// 	},
		// ),
		gobetterauthconfig.WithEventBus(
			gobetterauthmodels.EventBusConfig{
				Enabled: true,
				Prefix:  "gobetterauthplayground.",
				// Uncomment to test out using watermill with Kafka/Redpanda as the event bus.
				// The kafka publisher and subscriber are imported from this project's module.
				// PubSub: gobetterauthevents.NewWatermillPubSub(
				// 	events.NewKafkaPublisher(),
				// 	events.NewKafkaSubscriber(),
				// ),
			},
		),
		gobetterauthconfig.WithPlugins(
			gobetterauthmodels.PluginsConfig{
				Plugins: []gobetterauthmodels.Plugin{
					// Logger Plugin
					loggerplugin.NewLoggerPlugin(loggerplugin.LoggerPluginConfigOptions{
						MaxLogCount: 5,
					}),
					// Inline Plugin Example
					gobetterauthconfig.NewPlugin(
						gobetterauthconfig.WithPluginMetadata(
							gobetterauthmodels.PluginMetadata{
								Name:        "Inline Plugin",
								Version:     "0.0.1",
								Description: "This is an example of a simple inline plugin.",
							},
						),
						gobetterauthconfig.WithPluginConfig(
							gobetterauthmodels.PluginConfig{
								Enabled: true,
							},
						),
						gobetterauthconfig.WithPluginInit(func(ctx *gobetterauthmodels.PluginContext) error {
							slog.Info("Inline Plugin initialized for App: " + ctx.Config.AppName)
							return nil
						}),
						gobetterauthconfig.WithPluginRoutes(
							[]gobetterauthmodels.PluginRoute{
								{
									Method: "GET",
									Path:   "/inline/ping",
									Handler: func() http.Handler {
										return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
											w.Header().Set("Content-Type", "application/json")
											w.WriteHeader(200)
											json.NewEncoder(w).Encode(map[string]any{
												"message": "pong",
											})
										})
									},
								},
							},
						),
					),
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
			goBetterAuth.RedirectAuthMiddleware("https://go-better-auth.vercel.app", http.StatusSeeOther),
		},
		Handler: func(config *gobetterauthmodels.Config) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"message": fmt.Sprintf("%s: Hello from %s", config.AppName, req.URL.Path),
				})
			})
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
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"message": fmt.Sprintf("%s: data received on %s", config.AppName, r.URL.Path),
					"data":    body,
				})
			})
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
