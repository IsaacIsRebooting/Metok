package launcher

import (
	"context"

	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/components"
	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/components/mysqlx"
	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/gofer"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos/v2/config"
)

type ComponentsLauncher struct {
	group  *gofer.Group
	config map[string]config.Value
}

func NewComponentsLauncher(config config.Config) *ComponentsLauncher {
	configMap, err := config.Value("components").Map()
	if err != nil {
		panic("get components config failed, error: " + err.Error())
	}

	return &ComponentsLauncher{
		group: gofer.NewGroup(
			context.Background(),
			gofer.UserErrorGroup(),
		),
		config: configMap,
	}
}

func (l *ComponentsLauncher) Launch() {
	for componentsName, cfg := range l.config {
		log.Infof("launch component: %s", componentsName)
		log.Infof("%s config: %v", componentsName, cfg)
		launchWrapper(cfg, componentsName)
	}
	if err := l.group.Wait(); err != nil {
		log.Errorf("components launcher wait error: %v", err)
	}
}

func launchWrapper(cfg config.Value, componentsName string) {
	switch componentsName {
	case "mysql":
		launchComponent(cfg, mysqlx.Init)
	}
}
func launchComponent[T any](cfg config.Value, initMethod func(cfg components.ConfigMap[*T]) (func() error, error)) {
	configs, err := cfg.Map()
	if err != nil {
		panic("get component config failed, error: " + err.Error())
	}

	_, component := components.Load(configs, initMethod)
	if err := component.Start(); err != nil {
		panic("launch component failed, error: " + err.Error())
	}
}
