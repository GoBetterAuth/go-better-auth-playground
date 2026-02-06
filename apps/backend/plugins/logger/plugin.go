package logger

import (
	"context"
	"embed"
	"fmt"

	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/repositories"
	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/services"
	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/types"
	"github.com/GoBetterAuth/go-better-auth/v2/models"
	emailpasswordpluginconstants "github.com/GoBetterAuth/go-better-auth/v2/plugins/email-password/constants"
)

type LoggerPlugin struct {
	config        types.LoggerPluginConfig
	logger        models.Logger
	ctx           *models.PluginContext
	loggerService services.LoggerService
}

func New(config types.LoggerPluginConfig) *LoggerPlugin {
	return &LoggerPlugin{config: config}
}

func (p *LoggerPlugin) Metadata() models.PluginMetadata {
	return models.PluginMetadata{
		ID:          "logger",
		Version:     "1.0.0",
		Description: "Logs user authentication events to the database",
	}
}

func (p *LoggerPlugin) Config() any {
	return p.config
}

func (p *LoggerPlugin) Init(ctx *models.PluginContext) error {
	p.ctx = ctx
	p.logger = ctx.Logger

	if err := p.config.Validate(); err != nil {
		return fmt.Errorf("invalid logger plugin configuration: %w", err)
	}

	p.loggerService = services.NewService(repositories.NewBunLoggerRepository(ctx.DB), p.logger, p.config)

	p.subscribeToEvents()

	return nil
}

func (p *LoggerPlugin) Migrations(ctx context.Context, dbProvider string) (*embed.FS, error) {
	return GetMigrations(ctx, dbProvider)
}

func (p *LoggerPlugin) Routes() []models.Route {
	if p.ctx == nil || p.loggerService == nil {
		return nil
	}

	logger := p.ctx.Logger

	return Routes(logger, p.loggerService)
}

func (p *LoggerPlugin) Close() error {
	return nil
}

func (p *LoggerPlugin) subscribeToEvents() {
	_, err := p.ctx.EventBus.Subscribe(emailpasswordpluginconstants.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
		if _, err := p.loggerService.CreateLogEntry(ctx, emailpasswordpluginconstants.EventUserSignedUp, string(event.Payload)); err != nil {
			p.logger.Error("failed to create log entry for user sign up event", "error", err)
		}
		return nil
	})
	if err != nil {
		p.logger.Error("failed to subscribe to event", "event", emailpasswordpluginconstants.EventUserSignedUp, "error", err)
		return
	}
}
