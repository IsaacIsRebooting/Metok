package launcher

type App struct {
	Name          string `yaml:"name" json:"name"`
	Version       string `yaml:"version" json:"version"`
	Node          string `yaml:"node" json:"node"`
	TraceEndpoint string `yaml:"trace_endpoint" json:"trace_endpoint"`
}
