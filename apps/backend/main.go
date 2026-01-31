package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	gobetterauth "github.com/GoBetterAuth/go-better-auth/v2"
	gobetterauthconfig "github.com/GoBetterAuth/go-better-auth/v2/config"
	gobetterauthenv "github.com/GoBetterAuth/go-better-auth/v2/env"
	gobetterauthevents "github.com/GoBetterAuth/go-better-auth/v2/events"
	gobetterauthmodels "github.com/GoBetterAuth/go-better-auth/v2/models"

	bearerplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/bearer"
	configmanagerplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/config-manager"
	"github.com/GoBetterAuth/go-better-auth/v2/plugins/config-manager/types"
	csrfplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/csrf"
	emailplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/email"
	emailpasswordplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/email-password"
	emailpasswordplugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/email-password/types"
	emailplugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/email/types"
	jwtplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/jwt"
	jwtplugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/jwt/types"
	oauth2plugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/oauth2"
	oauth2plugintypes "github.com/GoBetterAuth/go-better-auth/v2/plugins/oauth2/types"
	ratelimitplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/rate-limit"
	secondarystorageplugin "github.com/GoBetterAuth/go-better-auth/v2/plugins/secondary-storage"

	loggerplugin "github.com/GoBetterAuth/go-better-auth-playground/plugins/logger"
	loggerplugintypes "github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/types"
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
					bearerplugin.HookIDBearerAuth.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/sign-in",
				Plugins: []string{
					bearerplugin.HookIDBearerAuthOptional.String(),
					jwtplugin.HookIDJWTRespondJSON.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/sign-up",
				Plugins: []string{
					bearerplugin.HookIDBearerAuthOptional.String(),
					jwtplugin.HookIDJWTRespondJSON.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/send-email-verification",
				Plugins: []string{
					bearerplugin.HookIDBearerAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/request-email-change",
				Plugins: []string{
					bearerplugin.HookIDBearerAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/sign-out",
				Plugins: []string{
					bearerplugin.HookIDBearerAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/tokens/refresh",
				Plugins: []string{
					bearerplugin.HookIDBearerAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Method: "GET",
				Path:   "/oauth2/callback/{provider}",
				Plugins: []string{
					bearerplugin.HookIDBearerAuthOptional.String(),
					jwtplugin.HookIDJWTRespondJSON.String(),
				},
			},
			{
				Method: "GET",
				Path:   "/api/protected",
				Plugins: []string{
					bearerplugin.HookIDBearerAuth.String(),
				},
			},
			{
				Method: "POST",
				Path:   "/api/protected",
				Plugins: []string{
					bearerplugin.HookIDBearerAuth.String(),
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
			configmanagerplugin.New(types.ConfigManagerPluginConfig{
				Enabled: false,
			}),
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
			// sessionplugin.New(sessionplugin.SessionPluginConfig{
			// 	Enabled: true,
			// }),
			jwtplugin.New(jwtplugintypes.JWTPluginConfig{
				Enabled:   true,
				Algorithm: jwtplugintypes.JWTAlgEdDSA,
				ExpiresIn: 30 * time.Second,
			}),
			bearerplugin.New(bearerplugin.BearerPluginConfig{
				Enabled: true,
			}),
			ratelimitplugin.New(ratelimitplugin.RateLimitPluginConfig{
				Enabled:  true,
				Provider: ratelimitplugin.RateLimitProviderRedis,
			}),
			loggerplugin.New(loggerplugintypes.LoggerPluginConfig{
				Enabled:     true,
				MaxLogCount: 10,
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
	})

	// goBetterAuth.RegisterHook(gobetterauthmodels.Hook{
	// 	Stage: gobetterauthmodels.HookAfter,
	// 	Matcher: func(ctx *gobetterauthmodels.RequestContext) bool {
	// 		return ctx.Method == "POST" && ctx.Path == "/api/auth/sign-up"
	// 	},
	// 	Handler: func(ctx *gobetterauthmodels.RequestContext) error {
	// 		ctx.ResponseHeaders.Set("Content-Type", "text/html")
	// 		ctx.ResponseStatus = http.StatusOK
	// 		ctx.ResponseBody = []byte("<html><body><h1>Sign Up Successful</h1></body></html>")
	// 		ctx.ResponseReady = true
	// 		ctx.Handled = true
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
