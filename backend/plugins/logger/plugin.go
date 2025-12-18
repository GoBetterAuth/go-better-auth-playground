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

type LogEntryDatabaseHooks struct {
	BeforeCreate func(entry *LogEntry) error
	AfterCreate  func(entry LogEntry) error
}

type LogEntryEventHooks struct {
	OnLogCreated func(entry LogEntry) error
}

type LoggerPluginConfigOptions struct {
	DatabaseHooks *LogEntryDatabaseHooks
	EventHooks    *LogEntryEventHooks
	MaxLogCount   int
}

type LoggerPlugin struct {
	// Embed the base plugin struct to inherit common functionality
	gobetterauthmodels.BasePlugin
	// service to handle business logic
	service *LoggerService
	// local state to keep track of the number of logs executed during runtime (just a random example, not stored in DB)
	logCount int
	// keep track of event subscriptions to allow for unsubscription on plugin close
	eventSubscriptions    map[string]gobetterauthmodels.SubscriptionID
	logEntryDatabaseHooks *LogEntryDatabaseHooks
	logEntryEventHooks    *LogEntryEventHooks
}

func NewLoggerPlugin(options LoggerPluginConfigOptions) gobetterauthmodels.Plugin {
	plugin := &LoggerPlugin{
		eventSubscriptions: make(map[string]gobetterauthmodels.SubscriptionID),
	}
	plugin.SetConfig(gobetterauthmodels.PluginConfig{
		Enabled: true,
		Options: options,
	})

	dbHooks := options.DatabaseHooks
	if dbHooks == nil {
		dbHooks = &LogEntryDatabaseHooks{}
	}
	plugin.SetDatabaseHooks(dbHooks)

	eventHooks := options.EventHooks
	if eventHooks == nil {
		eventHooks = &LogEntryEventHooks{}
	}
	plugin.SetEventHooks(eventHooks)

	return plugin
}

func (plugin *LoggerPlugin) Init(ctx *gobetterauthmodels.PluginContext) error {
	// Store the context for later use
	plugin.SetCtx(ctx)

	// Initialise your plugin's service (handles business logic for your plugin)
	plugin.service = NewLoggerService(ctx.Config.DB, plugin.logEntryDatabaseHooks)

	var eventUserSignedUpSubId gobetterauthmodels.SubscriptionID
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

				if user, err := plugin.Ctx().Api.Users.GetUserByID(userId.(string)); err != nil {
					slog.Error("failed to get user by ID", "error", err)
					return err
				} else {
					details := fmt.Sprintf("[LoggerPlugin] - User signed up: user_id=%s, email=%s", user.ID, user.Email)
					slog.Info(details)

					plugin.logCount++

					logEntry, err := plugin.service.CreateLogEntry(gobetterauthmodels.EventUserSignedUp, details)
					if err != nil {
						return err
					}

					if plugin.logEntryEventHooks != nil && plugin.logEntryEventHooks.OnLogCreated != nil {
						if err := plugin.logEntryEventHooks.OnLogCreated(*logEntry); err != nil {
							return err
						}
					}

					plugin.checkAndHandleMaxLogsReached(gobetterauthmodels.EventUserSignedUp, eventUserSignedUpSubId)
				}
			}

			return nil
		}); err != nil {
		return err
	} else {
		eventUserSignedUpSubId = id
		slog.Info("Subscribed to EventUserSignedUp", "subscription_id", eventUserSignedUpSubId)
		plugin.eventSubscriptions[gobetterauthmodels.EventUserSignedUp] = eventUserSignedUpSubId
	}

	var eventUserLoggedInSubId gobetterauthmodels.SubscriptionID
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

				if user, err := plugin.Ctx().Api.Users.GetUserByID(userId.(string)); err != nil {
					slog.Error("failed to get user by ID", "error", err)
					return err
				} else {
					details := fmt.Sprintf("[LoggerPlugin] - User logged in: user_id=%s, email=%s", user.ID, user.Email)
					slog.Info(details)

					plugin.logCount++

					logEntry, err := plugin.service.CreateLogEntry(gobetterauthmodels.EventUserLoggedIn, details)
					if err != nil {
						return err
					}

					if plugin.logEntryEventHooks != nil && plugin.logEntryEventHooks.OnLogCreated != nil {
						if err := plugin.logEntryEventHooks.OnLogCreated(*logEntry); err != nil {
							return err
						}
					}

					plugin.checkAndHandleMaxLogsReached(gobetterauthmodels.EventUserLoggedIn, eventUserLoggedInSubId)
				}
			}

			return nil
		}); err != nil {
		return err
	} else {
		eventUserLoggedInSubId = id
		slog.Info("Subscribed to EventUserLoggedIn", "subscription_id", eventUserLoggedInSubId)
		plugin.eventSubscriptions[gobetterauthmodels.EventUserLoggedIn] = eventUserLoggedInSubId
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
				plugin.Ctx().Middleware.Auth(),
			},
			Handler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]any{
						"logCount": plugin.logCount,
					})
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

func (plugin *LoggerPlugin) DatabaseHooks() any {
	return plugin.logEntryDatabaseHooks
}

func (plugin *LoggerPlugin) EventHooks() any {
	return plugin.logEntryEventHooks
}

func (plugin *LoggerPlugin) Close() error {
	for eventType, subId := range plugin.eventSubscriptions {
		plugin.Ctx().EventBus.Unsubscribe(eventType, subId)
	}
	slog.Info("LoggerPlugin closed and unsubscribed from events", "count", len(plugin.eventSubscriptions))

	return nil
}

// checkAndHandleMaxLogsReached checks if the maximum log count has been reached.
// If it has, it unsubscribes from the specified event type.
func (plugin *LoggerPlugin) checkAndHandleMaxLogsReached(eventType string, id gobetterauthmodels.SubscriptionID) {
	pluginOptions := plugin.Config().Options.(LoggerPluginConfigOptions)

	if plugin.logCount >= pluginOptions.MaxLogCount {
		plugin.Ctx().EventBus.Unsubscribe(eventType, id)
		slog.Info("LoggerPlugin unsubscribed from event due to max log count reached", "eventType", eventType, "maxLogCount", pluginOptions.MaxLogCount)
	}
}
