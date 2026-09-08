package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/Authula/authula"
	authulaconfig "github.com/Authula/authula/config"
	authulaenv "github.com/Authula/authula/env"
	authulaevents "github.com/Authula/authula/events"
	authulamodels "github.com/Authula/authula/models"
	accesscontrolplugin "github.com/Authula/authula/plugins/access-control"
	accesscontrolplugintypes "github.com/Authula/authula/plugins/access-control/types"
	adminplugin "github.com/Authula/authula/plugins/admin"
	adminplugintypes "github.com/Authula/authula/plugins/admin/types"
	csrfplugin "github.com/Authula/authula/plugins/csrf"
	emailplugin "github.com/Authula/authula/plugins/email"
	emailpasswordplugin "github.com/Authula/authula/plugins/email-password"
	emailpasswordplugintypes "github.com/Authula/authula/plugins/email-password/types"
	emailplugintypes "github.com/Authula/authula/plugins/email/types"
	magiclinkplugin "github.com/Authula/authula/plugins/magic-link"
	magiclinkplugintypes "github.com/Authula/authula/plugins/magic-link/types"
	oauth2plugin "github.com/Authula/authula/plugins/oauth2"
	oauth2plugintypes "github.com/Authula/authula/plugins/oauth2/types"
	organizationsplugin "github.com/Authula/authula/plugins/organizations"
	organizationsplugintypes "github.com/Authula/authula/plugins/organizations/types"
	ratelimitplugin "github.com/Authula/authula/plugins/rate-limit"
	ratelimitplugintypes "github.com/Authula/authula/plugins/rate-limit/types"
	secondarystorageplugin "github.com/Authula/authula/plugins/secondary-storage"
	sessionplugin "github.com/Authula/authula/plugins/session"

	bearerplugin "github.com/Authula/authula/plugins/bearer"
	jwtplugin "github.com/Authula/authula/plugins/jwt"
	jwtplugintypes "github.com/Authula/authula/plugins/jwt/types"

	loggerplugin "github.com/Authula/authula-playground/plugins/logger"
	loggerplugintypes "github.com/Authula/authula-playground/plugins/logger/types"
)

type EmailPasswordServiceHooks = emailpasswordplugintypes.EmailPasswordServiceHooksConfig

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
	// Init Authula config
	// -------------------------------------

	config := authulaconfig.NewConfig(
		authulaconfig.WithAppName("AuthulaPlayground"),
		authulaconfig.WithBasePath("/api/auth"),
		authulaconfig.WithDatabase(authulamodels.DatabaseConfig{
			Provider: "postgres",
			URL:      os.Getenv(authulaenv.EnvDatabaseURL),
		}),
		authulaconfig.WithLogger(authulamodels.LoggerConfig{
			Level: "debug",
		}),
		authulaconfig.WithSession(authulamodels.SessionConfig{
			CookieName:         "authula.session_token",
			ExpiresIn:          24 * time.Hour,
			UpdateAge:          5 * time.Minute,
			CookieMaxAge:       24 * time.Hour,
			Secure:             false,
			HttpOnly:           true,
			SameSite:           "lax",
			MaxSessionsPerUser: 5,
			AutoCleanup:        true,
			CleanupInterval:    time.Minute,
		}),
		authulaconfig.WithVerification(authulamodels.VerificationConfig{
			AutoCleanup:     true,
			CleanupInterval: time.Minute,
		}),
		authulaconfig.WithSecurity(authulamodels.SecurityConfig{
			TrustedOrigins: []string{"http://localhost:3000"},
			CORS: authulamodels.CORSConfig{
				AllowCredentials: true,
				AllowedOrigins:   []string{"http://localhost:3000"},
				AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PATCH", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Authorization", "Content-Type", "Set-Cookie", "Cookie", "X-AUTHULA-CSRF-TOKEN"},
				ExposedHeaders:   []string{"X-AUTHULA-CSRF-TOKEN"},
				MaxAge:           24 * time.Hour,
			},
		}),
		authulaconfig.WithEventBus(authulamodels.EventBusConfig{
			Provider: authulaevents.ProviderKafka,
			Kafka: &authulamodels.KafkaConfig{
				Brokers:       os.Getenv(authulaenv.EnvKafkaBrokers),
				ConsumerGroup: os.Getenv(authulaenv.EnvEventBusConsumerGroup),
			},
		}),
		authulaconfig.WithRouteMappings([]authulamodels.RouteMapping{
			// Core Routes
			{
				Paths: []string{"GET:/me"},
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					// bearer.HookIDBearerAuth.String(),
					// jwt.HookIDJWTRespondJSON.String(),
				},
			},
			{
				Paths: []string{"POST:/sign-out"},
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					// bearer.HookIDBearerAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			// Email-Password Routes
			{
				Paths: []string{
					"POST:/email-password/sign-in",
					"POST:/email-password/sign-up",
					"POST:/email-password/request-password-reset",
					"POST:/email-password/change-password",
				},
				Plugins: []string{
					sessionplugin.HookIDSessionAuthOptional.String(),
					// bearer.HookIDBearerAuthOptional.String(),
					// jwt.HookIDJWTRespondJSON.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			{
				Paths: []string{"GET:/email-password/verify-email"},
				Plugins: []string{
					sessionplugin.HookIDSessionAuthOptional.String(),
					// bearer.HookIDBearerAuthOptional.String(),
				},
			},
			{
				Paths: []string{
					"POST:/email-password/send-email-verification",
					"POST:/email-password/request-email-change",
				},
				Plugins: []string{
					sessionplugin.HookIDSessionAuth.String(),
					// bearer.HookIDBearerAuth.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			// Magic Link Routes
			{
				Paths: []string{
					"POST:/magic-link/sign-in",
					"POST:/magic-link/verify",
					"POST:/magic-link/exchange",
				},
				Plugins: []string{
					sessionplugin.HookIDSessionAuthOptional.String(),
					// bearer.HookIDBearerAuthOptional.String(),
					// jwt.HookIDJWTRespondJSON.String(),
					csrfplugin.HookIDCSRFProtect.String(),
				},
			},
			// Custom Routes
			{
				Paths:   []string{"GET:/api/health"},
				Plugins: []string{},
			},
		}),
	)

	// -------------------------------------
	// Init Authula instance
	// -------------------------------------

	auth := authula.New(&authula.AuthConfig{
		Config: config,
		Plugins: []authulamodels.Plugin{
			// Built-in plugins
			// Secondary storage plugin MUST be registered before rate-limit plugin
			// This allows rate-limit to optionally use Redis/database for distributed rate limiting
			secondarystorageplugin.New(secondarystorageplugin.SecondaryStoragePluginConfig{
				Enabled:  true,
				Provider: secondarystorageplugin.SecondaryStorageProviderRedis,
				Redis: &secondarystorageplugin.RedisStorageConfig{
					URL:         os.Getenv(authulaenv.EnvRedisURL),
					MaxRetries:  3,
					PoolSize:    10,
					PoolTimeout: 30 * time.Second,
				},
			}),
			csrfplugin.New(csrfplugin.CSRFPluginConfig{
				Enabled: true,
			}),
			emailplugin.New(emailplugintypes.EmailPluginConfig{
				Enabled:     true,
				Provider:    emailplugintypes.ProviderSMTP,
				FromAddress: "noreply@example.com",
				TLSMode:     emailplugintypes.SMTPTLSModeStartTLS,
			}),
			emailpasswordplugin.New(emailpasswordplugintypes.EmailPasswordPluginConfig{
				Enabled:                     true,
				MinPasswordLength:           8,
				MaxPasswordLength:           32,
				DisableSignUp:               false,
				RequireEmailVerification:    true,
				AutoSignIn:                  true,
				SendEmailOnSignUp:           true,
				SendEmailOnSignIn:           false,
				EmailVerificationExpiresIn:  24 * time.Hour,
				PasswordResetExpiresIn:      time.Hour,
				RequestEmailChangeExpiresIn: time.Hour,
			}),
			oauth2plugin.New(oauth2plugintypes.OAuth2PluginConfig{
				Enabled: true,
				Providers: map[string]oauth2plugintypes.ProviderConfig{
					"discord": {
						Enabled:      true,
						ClientID:     os.Getenv(authulaenv.EnvDiscordClientID),
						ClientSecret: os.Getenv(authulaenv.EnvDiscordClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/discord", config.BaseURL, config.BasePath),
					},
					"github": {
						Enabled:      true,
						ClientID:     os.Getenv(authulaenv.EnvGithubClientID),
						ClientSecret: os.Getenv(authulaenv.EnvGithubClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/github", config.BaseURL, config.BasePath),
					},
					"google": {
						Enabled:      true,
						ClientID:     os.Getenv(authulaenv.EnvGoogleClientID),
						ClientSecret: os.Getenv(authulaenv.EnvGoogleClientSecret),
						RedirectURL:  fmt.Sprintf("%s%s/oauth2/callback/google", config.BaseURL, config.BasePath),
					},
				},
			}),
			sessionplugin.New(sessionplugin.SessionPluginConfig{
				Enabled: true,
			}),
			jwtplugin.New(jwtplugintypes.JWTPluginConfig{
				Enabled: false,
			}),
			bearerplugin.New(bearerplugin.BearerPluginConfig{
				Enabled: false,
			}),
			magiclinkplugin.New(magiclinkplugintypes.MagicLinkPluginConfig{
				Enabled:       true,
				ExpiresIn:     time.Hour,
				DisableSignUp: false,
			}),
			accesscontrolplugin.New(accesscontrolplugintypes.AccessControlPluginConfig{
				Enabled: true,
			}),
			adminplugin.New(adminplugintypes.AdminPluginConfig{
				Enabled:                   true,
				ImpersonationMaxExpiresIn: 15 * time.Minute,
			}),
			organizationsplugin.New(organizationsplugintypes.OrganizationsPluginConfig{
				Enabled:                          true,
				RequireEmailVerifiedOnInvitation: true,
			}),
			ratelimitplugin.New(ratelimitplugintypes.RateLimitPluginConfig{
				Enabled:     true,
				Provider:    ratelimitplugintypes.RateLimitProviderRedis,
				Window:      time.Minute,
				Max:         60,
				CustomRules: map[string]ratelimitplugintypes.RateLimitRule{},
			}),

			// Custom plugins
			loggerplugin.New(loggerplugintypes.LoggerPluginConfig{
				Enabled:     true,
				MaxLogCount: 10,
			}),
		},
	})

	// You can uncomment the following 2 lines to drop all migrations (i.e., reset the database).
	// ctx := context.Background()
	// if err := authula.PluginRegistry.DropMigrations(ctx); err != nil {
	// 	slog.Error("failed to drop plugin migrations", "error", err)
	// 	return
	// }
	// if err := authula.DropCoreMigrations(ctx); err != nil {
	// 	slog.Error("failed to drop core migrations", "error", err)
	// 	return
	// }
	// slog.Info("all migrations dropped successfully")
	// return

	// // Example of how to programmatically interact with plugins outside of route handlers, hooks, etc.
	// organizationsPlugin, ok := authula.PluginRegistry.GetPlugin(authulamodels.PluginOrganizations.String()).(*organizations.OrganizationsPlugin)
	// if !ok {
	// 	// You can return early here in your own code or handle it differently however you wish...
	// 	return
	// }

	// // Now you can use the plugin's API to perform actions programmatically outside of route handlers, hooks, etc.
	// ctx := context.Background()
	// organizations, err := organizationsPlugin.Api.GetAllOrganizations(ctx, "f8863909-2e90-4276-bab6-84659550b294")
	// if err != nil {
	// 	slog.Error("failed to get organizations", "error", err)
	// } else {
	// 	slog.Info("organizations retrieved successfully", "organizations", organizations)
	// }

	// ctx := context.Background()
	// signOutResult, err := auth.Api.SignOut(ctx, "2c287cae-70bb-4a97-b4a8-313e301bca7b", nil, new(true))
	// if err != nil {
	// 	slog.Error("failed to sign out", "error", err)
	// } else {
	// 	slog.Info(signOutResult.Message)
	// }

	// -------------------------------------
	// Add custom routes to the router
	// Note: Call RegisterCustomRoute() BEFORE Handler() to ensure routes are registered before handler is served
	// Custom routes are registered without the /api/auth prefix
	// -------------------------------------

	// Health check endpoint
	auth.RegisterCustomRoute(authulamodels.Route{
		Method: "GET",
		Path:   "/api/health",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx, _ := authulamodels.GetRequestContext(r.Context())
			reqCtx.SetJSONResponse(http.StatusOK, map[string]any{
				"status": "ok",
			})
		}),
	})

	// authula.RegisterHook(authulamodels.Hook{
	// 	Stage: authulamodels.HookBefore,
	// 	Matcher: func(ctx *authulamodels.RequestContext) bool {
	// 		return ctx.UserID != nil && *ctx.UserID != "" && slices.Contains(
	// 			[]string{
	// 				"/api/protected",
	// 				"/path/to/more/routes...",
	// 			},
	// 			ctx.Path,
	// 		)
	// 	},
	// 	Handler: func(ctx *authulamodels.RequestContext) error {
	// 		// Do as you wish before the request is processed by the route handler...
	// 		return nil
	// 	},
	// })

	// -------------------------------------
	// Attach Authula handler to your chosen framework and run your server
	// All hooks (CORS, auth, rate limiting, etc.) are applied via the plugin system
	// -------------------------------------

	port := os.Getenv(authulaenv.EnvPort)
	slog.Debug(fmt.Sprintf("Server running on http://localhost:%s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), auth.Handler()); err != nil {
		slog.Error("Server error", "err", err)
	}
}
