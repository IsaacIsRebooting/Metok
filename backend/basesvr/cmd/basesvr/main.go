package main

import (
	"context"

	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/conf"
	"github.com/IsaacIsRebooting/Metok/backend/basesvr/internal/data/utils"
	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/components/mysqlx"
	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/launcher"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
// var (
// 	// Name is the name of the compiled software.
// 	Name string
// 	// Version is the version of the compiled software.
// 	Version string
// 	// flagconf is the config flag.
// 	flagconf string

// 	id, _ = os.Hostname()
// )

// func init() {
// 	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
// }

// func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
// 	return kratos.New(
// 		kratos.ID(id),
// 		kratos.Name(Name),
// 		kratos.Version(Version),
// 		kratos.Metadata(map[string]string{}),
// 		kratos.Logger(logger),
// 		kratos.Server(
// 			gs,
// 			hs,
// 		),
// 	)
// }

// func main() {
// 	flag.Parse()
// 	logger := log.With(log.NewStdLogger(os.Stdout),
// 		"ts", log.DefaultTimestamp,
// 		"caller", log.DefaultCaller,
// 		"service.id", id,
// 		"service.name", Name,
// 		"service.version", Version,
// 		"trace.id", tracing.TraceID(),
// 		"span.id", tracing.SpanID(),
// 	)
// 	c := config.New(
// 		config.WithSource(
// 			file.NewSource(flagconf),
// 		),
// 	)
// 	defer c.Close()

// 	if err := c.Load(); err != nil {
// 		panic(err)
// 	}

// 	var bc conf.Bootstrap
// 	if err := c.Scan(&bc); err != nil {
// 		panic(err)
// 	}

// 	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer cleanup()

// 	// start and wait for stop signal
// 	if err := app.Run(); err != nil {
// 		panic(err)
// 	}
// }

func main() {
	c := &conf.Config{}
	launcher.New(
		launcher.WithConfigValue(c),
		launcher.WithConfigOptions(
			config.WithSource(file.NewSource("configs/")),
		),
		launcher.WithAfterServerStartHandler(func() {
			query.SetDefault(mysqlx.GetDBClient(context.Background()))
		}),
		launcher.WithGrpcServer(func(configValue interface{}) *grpc.Server {
			cfg, ok := configValue.(*conf.Config)
			if !ok {
				panic("invalid config value")
			}

			utils.InitDefaultSnowflakeNode(cfg.Snowflake.Node)

			return server.NewGRPCServer(
				server.WithFileTableShardingConfig(cfg.Data),
				server.WithDBShardingTablesConfig(cfg.Data.DbShardingTables),
			)
		}),
	).Run()
}
