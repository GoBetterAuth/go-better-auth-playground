package logger

import (
	"net/http"

	"github.com/GoBetterAuth/go-better-auth-playground/plugins/logger/services"
	"github.com/GoBetterAuth/go-better-auth/v2/models"
)

const rateLimitKey = "plugin:logger:count"

// Routes creates and returns the plugin routes
func Routes(logger models.Logger, service services.LoggerService) []models.Route {
	logCountHandler := &LogCountHandler{
		service: service,
		logger:  logger,
	}

	return []models.Route{
		{
			Method:  http.MethodGet,
			Path:    "/logger/count",
			Handler: logCountHandler.Handler(),
		},
	}
}

type LogCountHandler struct {
	service services.LoggerService
	logger  models.Logger
}

func (h *LogCountHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqCtx, _ := models.GetRequestContext(ctx)

		if r.Method != http.MethodGet {
			reqCtx.SetJSONResponse(http.StatusMethodNotAllowed, map[string]any{
				"message": "only GET method allowed",
			})
			reqCtx.Handled = true
			return
		}

		logCount, err := h.service.GetLogCount(r.Context())
		if err != nil {
			reqCtx.SetJSONResponse(http.StatusInternalServerError, map[string]any{
				"message": "failed to get log count",
			})
			reqCtx.Handled = true
			return
		}

		reqCtx.SetJSONResponse(http.StatusOK, map[string]any{
			"logCount": logCount,
		})
	}
}
