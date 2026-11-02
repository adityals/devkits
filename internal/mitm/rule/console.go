package rule

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/dop251/goja"
)

var consoleLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"log":   slog.LevelInfo,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// installConsole exposes console.* in the VM, writing to logger tagged with the rule id.
func installConsole(vm *goja.Runtime, logger *slog.Logger, id string) error {
	console := vm.NewObject()
	for name, level := range consoleLevels {
		err := console.Set(name, func(call goja.FunctionCall) goja.Value {
			logger.Log(context.Background(), level, formatArgs(call.Arguments), "rule", id)
			return goja.Undefined()
		})
		if err != nil {
			return err
		}
	}

	return vm.Set("console", console)
}

func formatArgs(args []goja.Value) string {
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, formatValue(a))
	}

	return strings.Join(parts, " ")
}

// formatValue prints objects as JSON like a browser console; errors and functions use their string form.
func formatValue(v goja.Value) string {
	obj, ok := v.(*goja.Object)
	if !ok || obj.ClassName() == "Error" || obj.ClassName() == "Function" {
		return v.String()
	}

	b, err := json.Marshal(obj.Export())
	if err != nil {
		return v.String()
	}

	return string(b)
}
