package plugin

import (
	"fmt"
	"os"
	"plugin"
	"strings"
)

type HandlerPlugin interface {
	Handle(payload any) any
}

func Load(pluginDir string) (map[string]*plugin.Plugin, error) {
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil, err
	}

	plugins := make(map[string]*plugin.Plugin, len(files))

	for i := 0; i < len(files); i++ {
		plg, err := plugin.Open(fmt.Sprintf("%s/%s", pluginDir, files[i].Name()))
		if err != nil {
			return nil, err
		}

		fmt.Println(fmt.Sprintf("%s/%s", pluginDir, files[i].Name()))

		pluginName := strings.TrimSuffix(files[i].Name(), ".so")

		plugins[pluginName] = plg
	}

	return plugins, nil
}
