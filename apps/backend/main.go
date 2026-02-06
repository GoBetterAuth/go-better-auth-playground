package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	gobetterauth "github.com/GoBetterAuth/go-better-auth/v2"
	gobetterauthconfig "github.com/GoBetterAuth/go-better-auth/v2/config"
	gobetterauthenv "github.com/GoBetterAuth/go-better-auth/v2/env"
	gobetterauthevents "github.com/GoBetterAuth/go-better-auth/v2/events"
	gobetterauthmodels "github.com/GoBetterAuth/go-better-auth/v2/models"
	csrfplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/csrf"
	emailplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/email"
	emailpasswordplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/email-password"
	emailpasswordplugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/email-password/types"
	emailplugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/email/types"
	oauth2plugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/oauth2"
	oauth2plugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/oauth2/types"
	ratelimitplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/rate-limit"
	secondarystorageplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/secondary-storage"
	sessionplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/session"
)

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
			Provider: "postgres",
			URL:      os.Getenv(gobetterauthenv.EnvDatabaseURL),
		}),
		gobetterauthconfig.WithLogger(gobetterauthmodels.LoggerConfig{
			Level: "debug",
		}),
		gobetterauthconfig.WithSecurity(gobetterauthmodels.SecurityConfig{
			TrustedOrigins: []string{"http://localhost:3000"},
			CORS: gobetterauthmodels.CORSConfig{
				AllowCredentials: true,
				AllowedOrigins:   []string{"http://localhost:3000"},
				AllowedHeaders:   []string{"Authorization", "Content-Type", "Cookie", "Set-Cookie", "X-GOBETTERAUTH-CSRF-TOKEN"},
				ExposedHeaders:   []string{"X-GOBETTERAUTH-CSRF-TOKEN"},
			},
		}),
		gobetterauthconfig.WithEventBus(gobetterauthmodels.EventBusConfig{
			Provider: gobetterauthevents.ProviderKafka,
			Kafka: &gobetterauthmodels.KafkaConfig{
				Brokers:       os.Getenv(gobetterauthenv.EnvKafkaBrokers),
				ConsumerGroup: os.Getenv(gobetterauthenv.EnvEventBusConsumerGroup),
			},
		}),
		gobetterauthconfig.WithRouteMappings([]gobetterauthmodels.RouteMapping{
			{
				Method: "GET",
				Path:   "/me",
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/sign-in",
				Plugins: []string{
					sessionplugin.HookIDSessionAuthOptional.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/sign-up",
				Plugins: []string{
					sessionplugin.HookIDSessionAuthOptional.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/send-email-verification",
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/request-email-change",
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/sign-out",
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/tokens/refresh",
				Plugins: []string{
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			// {
			// 	Method:  "GET",
			// 	Path:    "/oauth2/callback/{provider}",
			// 	Plugins: []string{
			// 		// bearerplugin.HookIDBearerAuthOptional.String(),
			// 		// jwtplugin.HookIDJWTRespondJSON.String(),
			// 	},
			// },
			{
				Method: "GET",
				Path:   "/api/protected",
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/api/protected",
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
		}),
	)

	// -------------------------------------
	// Init GoBetterAuth instance
	// -------------------------------------

	goBetterAuth := gobetterauth.New(&gobetterauth.AuthConfig{
		Config: config,
		Plugins: []gobetterauthmodels.Plugin{
			// Secondary storage plugin MUST be registered before rate-limit plugin
			// This allows rate-limit to optionally use Redis/database for distributed rate limiting
			secondarystorageplugin.New(secondarystorageplugin.SecondaryStoragePluginConfig{
				Enabled:  true,
				Provider: secondarystorageplugin.SecondaryStorageProviderRedis,
				Redis: &secondarystorageplugin.RedisStorageConfig{
					URL: os.Getenv(gobetterauthenv.EnvRedisURL),
				},
			}),
			csrfplugin.New(csrfplugin.CSRFPluginConfig{
				Enabled: true,
			}),
			emailplugin.New(emailplugintypes.EmailPluginConfig{
				Enabled:          true,
				Provider:         emailplugintypes.ProviderSMTP,
				FallbackProvider: emailplugintypes.ProviderResend,
				FromAddress:      "from@example.com",
			}),
			emailpasswordplugin.New(emailpasswordplugintypes.EmailPasswordPluginConfig{
				Enabled:                  true,
				MinPasswordLength:        8,
				MaxPasswordLength:        32,
				DisableSignUp:            false,
				RequireEmailVerification: true,
				AutoSignIn:               true,
				SendEmailOnSignUp:        true,
			}),
			oauth2plugin.New(oauth2plugintypes.OAuth2PluginConfig{
				Enabled: true,
				Providers: map[string]oauth2plugintypes.ProviderConfig{
					"discord": {
						Enabled:      true,
						ClientID:     os.Getenv(gobetterauthenv.EnvDiscordClientID),
						ClientSecret: os.Getenv(gobetterauthenv.EnvDiscordClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/discord", config.BaseURL, config.BasePath),
					},
					"github": {
						Enabled:      true,
						ClientID:     os.Getenv(gobetterauthenv.EnvGithubClientID),
						ClientSecret: os.Getenv(gobetterauthenv.EnvGithubClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/github", config.BaseURL, config.BasePath),
					},
					"google": {
						Enabled:      true,
						ClientID:     os.Getenv(gobetterauthenv.EnvGoogleClientID),
						ClientSecret: os.Getenv(gobetterauthenv.EnvGoogleClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/google", config.BaseURL, config.BasePath),
					},
				},
			}),
			sessionplugin.New(sessionplugin.SessionPluginConfig{
				Enabled: true,
			}),
			ratelimitplugin.New(ratelimitplugin.RateLimitPluginConfig{
				Enabled:  true,
				Provider: ratelimitplugin.RateLimitProviderRedis,
			}),
		},
	})

	// You can uncomment the following 2 lines to drop all migrations (i.e., reset the database).
	// ctx := context.Background()
	// goBetterAuth.PluginRegistry.DropMigrations(ctx)
	// goBetterAuth.DropCoreMigrations(ctx)
	// return

	// -------------------------------------
	// Add custom routes to the router
	// Note: Call RegisterCustomRoute() BEFORE Handler() to ensure routes are registered before handler is served
	// Custom routes are registered without the /api/auth prefix
	// -------------------------------------

	// Health check endpoint
	goBetterAuth.RegisterCustomRoute(gobetterauthmodels.Route{
		Method:   "GET",
		Path:     "/api/health",
		Metadata: map[string]any{},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx, _ := gobetterauthmodels.GetRequestContext(r.Context())
			reqCtx.SetJSONResponse(http.StatusOK, map[string]any{
				"status": "ok",
			})
		}),
	})

	// Protected test endpoint
	goBetterAuth.RegisterCustomRoute(gobetterauthmodels.Route{
		Method: "GET",
		Path:   "/api/protected",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, _ := gobetterauthmodels.GetUserIDFromContext(r.Context())
			json.NewEncoder(w).Encode(map[string]any{
				"message": fmt.Sprintf("Hello, your user ID is %s", userId),
			})
		}),
		Metadata: map[string]any{
			"plugins": []string{sessionplugin.HookIDSessionAuth.String()},
		},
	})

	goBetterAuth.RegisterCustomRoute(gobetterauthmodels.Route{
		Method: "POST",
		Path:   "/api/protected",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, _ := gobetterauthmodels.GetUserIDFromContext(r.Context())
			json.NewEncoder(w).Encode(map[string]any{
				"message": fmt.Sprintf("Hello, your user ID is %s", userId),
			})
		}),
		Metadata: map[string]any{
			"plugins": []string{sessionplugin.HookIDSessionAuth.String()},
		},
	})

	// goBetterAuth.RegisterHook(gobetterauthmodels.Hook{
	// 	Stage: gobetterauthmodels.HookBefore,
	// 	Matcher: func(ctx *gobetterauthmodels.RequestContext) bool {
	// 		return ctx.UserID != nil && *ctx.UserID != "" && slices.Contains(
	// 			[]string{
	// 				"/api/protected",
	// 				"/path/to/more/routes...",
	// 			},
	// 			ctx.Path,
	// 		)
	// 	},
	// 	Handler: func(ctx *gobetterauthmodels.RequestContext) error {
	// 		// Do as you wish before the request is processed by the route handler...
	// 		return nil
	// 	},
	// })

	// -------------------------------------
	// Attach GoBetterAuth handler to your chosen framework and run your server
	// All hooks (CORS, auth, rate limiting, etc.) are applied via the plugin system
	// -------------------------------------

	port := os.Getenv(gobetterauthenv.EnvPort)
	slog.Debug(fmt.Sprintf("Server running on http://localhost:%s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), goBetterAuth.Handler()); err != nil {
		slog.Error("Server error", "err", err)
	}
}
