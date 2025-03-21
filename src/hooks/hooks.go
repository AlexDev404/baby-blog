package hooks

import (
	"baby-blog/types"
	"log/slog"
	"os"
)

type HooksConnector struct {
	Logger *slog.Logger
}

func Hooks(pageData map[string]interface{}, dbModels *types.Models) map[string]interface{} {
	// Create a handler that prepends "[HOOKS]" to log messages
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)

	// Create the hooks connector with the custom logger
	hooks := HooksConnector{
		Logger: slog.New(handler).With(slog.String("component", "Hooks")),
	}

	// Call all the hooks that are needed
	if pageData["Path"] == "feedback/gallery" {
		pageData = hooks.PageLoad(pageData, dbModels)
	}

	return pageData
}
