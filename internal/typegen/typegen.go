package typegen

import (
	"devkits/internal/typegen/convert"
	"devkits/internal/typegen/gots"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Services wires the typegen module. Each conversion route is one entry in the map.
func Services() []application.Service {
	routes := map[convert.Route]convert.Converter{
		{From: convert.LangGo, To: convert.LangTS}: gots.New(),
	}

	return []application.Service{
		application.NewService(convert.NewConvertService(routes)),
	}
}
