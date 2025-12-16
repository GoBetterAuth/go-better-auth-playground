package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	gobetterauthmodels "github.com/GoBetterAuth/go-better-auth/models"
)

// LoggerPlugin is a sample plugin used for demonstration purposes that carries out logging upon user events and stores log entries in the database.
type LoggerPlugin struct {
	config  gobetterauthmodels.PluginConfig
	ctx     *gobetterauthmodels.PluginContext
	service *LoggerService
	// local state to keep track of the number of logs executed during runtime (just a random example, not stored in DB)
	logCount int
	// keep track of event subscriptions to allow for unsubscription on plugin close
	eventSubscriptions map[string]gobetterauthmodels.SubscriptionID
}

func NewLoggerPlugin(config gobetterauthmodels.PluginConfig) *LoggerPlugin {
	return &LoggerPlugin{
		config:             config,
		logCount:           0,
		eventSubscriptions: make(map[string]gobetterauthmodels.SubscriptionID),
	}
}

func (plugin *LoggerPlugin) Metadata() gobetterauthmodels.PluginMetadata {
	return gobetterauthmodels.PluginMetadata{
		Name:        "Logger Plugin",
		Version:     "0.0.1",
		Description: "Logger plugin example",
	}
}

func (plugin *LoggerPlugin) Config() gobetterauthmodels.PluginConfig {
	return plugin.config
}

func (plugin *LoggerPlugin) Ctx() *gobetterauthmodels.PluginContext {
	return plugin.ctx
}

func (plugin *LoggerPlugin) Init(ctx *gobetterauthmodels.PluginContext) error {
	// Store the context for later use
	plugin.ctx = ctx

	// Initialise your plugin's service (handles business logic for your plugin)
	plugin.service = NewLoggerService(ctx.Config.DB)

	if id, err := ctx.EventBus.Subscribe(
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

				if user, err := plugin.ctx.Api.Users.GetUserByID(userId.(string)); err != nil {
					slog.Error("failed to get user by ID", "error", err)
					return err
				} else {
					details := fmt.Sprintf("[LoggerPlugin] - User signed up: user_id=%s, email=%s", user.ID, user.Email)
					slog.Info(details)
					plugin.logCount++
					return plugin.service.Log(gobetterauthmodels.EventUserSignedUp, details)
				}
			}

			return nil
		}); err != nil {
		return err
	} else {
		slog.Info("Subscribed to EventUserSignedUp", "subscription_id", id)
		plugin.eventSubscriptions[gobetterauthmodels.EventUserSignedUp] = id
	}

	if id, err := ctx.EventBus.Subscribe(
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

				if user, err := plugin.ctx.Api.Users.GetUserByID(userId.(string)); err != nil {
					slog.Error("failed to get user by ID", "error", err)
					return err
				} else {
					details := fmt.Sprintf("[LoggerPlugin] - User logged in: user_id=%s, email=%s", user.ID, user.Email)
					slog.Info(details)
					plugin.logCount++
					return plugin.service.Log(gobetterauthmodels.EventUserLoggedIn, details)
				}
			}

			return nil
		}); err != nil {
		return err
	} else {
		slog.Info("Subscribed to EventUserLoggedIn", "subscription_id", id)
		plugin.eventSubscriptions[gobetterauthmodels.EventUserLoggedIn] = id
	}

	return nil
}

func (plugin *LoggerPlugin) Migrations() []any {
	return []any{&LogEntry{}}
}

func (plugin *LoggerPlugin) Routes() []gobetterauthmodels.PluginRoute {
	return []gobetterauthmodels.PluginRoute{
		{
			Method: "GET",
			Path:   "/logger/count",
			Middleware: []gobetterauthmodels.PluginRouteMiddleware{
				plugin.ctx.Middleware.Auth(),
			},
			Handler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					jsonData := map[string]any{
						"logCount": plugin.logCount,
					}
					response, _ := json.Marshal(jsonData)
					w.Header().Set("Content-Type", "application/json")
					w.Write(response)
				})
			},
		},
	}
}

func (plugin *LoggerPlugin) RateLimit() *gobetterauthmodels.PluginRateLimit {
	return &gobetterauthmodels.RateLimitConfig{
		Enabled: true,
		Window:  30 * time.Second,
		Max:     1,
	}
}

func (plugin *LoggerPlugin) DatabaseHooks() *gobetterauthmodels.PluginDatabaseHooks {
	return nil
}

func (plugin *LoggerPlugin) EventHooks() *gobetterauthmodels.PluginEventHooks {
	return nil
}

func (plugin *LoggerPlugin) Close() error {
	for eventType, subId := range plugin.eventSubscriptions {
		plugin.ctx.EventBus.Unsubscribe(eventType, subId)
	}
	slog.Info("LoggerPlugin closed and unsubscribed from events", "count", len(plugin.eventSubscriptions))

	return nil
}
