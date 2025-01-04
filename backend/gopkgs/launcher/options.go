package launcher

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2/config"
)

type Option func(l *Launcher)

func WithBeforeConfigInitHandler(handler func()) Option {
	return func(l *Launcher) {
		l.beforeConfigInitHandlers = append(l.beforeConfigInitHandlers, handler)
	}
}

func WithAfterConfigInitHandler(handler func()) Option {
	return func(l *Launcher) {
		l.afterConfigInitHandlers = append(l.afterConfigInitHandlers, handler)
	}
}
func WithBeforeServerStartHandler(handler func()) Option {
	return func(l *Launcher) {
		l.beforeServerStartHandlers = append(l.beforeServerStartHandlers, handler)
	}
}

// WithAfterServerStartHandler 是一个用于注册服务器启动后处理函数的辅助函数。
// 它允许用户在服务器启动后执行自定义的初始化操作。
//
// 参数:
//
//	handler - 在服务器启动后要执行的处理函数。
//
// 返回值:
//
//	Option - 一个配置选项，用于在Launcher初始化时应用。
func WithAfterServerStartHandler(handler func()) Option {
	return func(l *Launcher) {
		l.afterServerStartHandlers = append(l.afterServerStartHandlers, handler)
	}
}

func WithShutdownHandler(handler func()) Option {
	return func(l *Launcher) {
		l.shutdownHandlers = append(l.shutdownHandlers, handler)
	}
}

// WithConfigValue 是一个函数，用于创建一个Option类型的函数，
// 这个函数当被调用时，会将一个配置值赋给Launcher类型的指针。
// 此设计模式允许使用者在创建Launcher实例后，通过传递不同的配置值来配置实例。
//
// 参数:
//
//	configValue - 类型为interface{}，代表可以是任何类型的配置值。
//
// 返回值:
//
//	返回一个Option类型的函数，该函数接受一个Launcher类型的指针。
//	这个返回的函数当被调用时，会将Launcher实例的configValue字段设置为传入的configValue值。
func WithConfigValue(configValue interface{}) Option {
	return func(l *Launcher) {
		l.configValue = configValue
	}
}

// WithConfigOptions 是一个返回 Option 类型的函数，用于将配置选项添加到 Launcher 的配置选项列表中。
// 参数 options 是一个变长参数，类型为 config.Option，代表一个或多个配置选项。
// 该函数返回一个函数类型，该内部函数接受一个 Launcher 类型的指针作为参数。
// 当内部函数被调用时，它会将传入的配置选项添加到 Launcher 实例的配置选项列表中。
func WithConfigOptions(options ...config.Option) Option {
	return func(l *Launcher) {
		l.configOptions = append(l.configOptions, options...)
	}
}

func WithConfigWatcher(key string, observer config.Observer) Option {
	return func(l *Launcher) {
		l.configWatchMap[key] = observer
	}
}
func WithHttpServer(s func(configValue interface{}) *http.Server) Option {
	return func(l *Launcher) {
		l.ginServer = s
	}
}
func WithGrpcServer(s func(configValue interface{}) *grpc.Server) Option {
	return func(l *Launcher) {
		l.grpcServer = s
	}
}
func WithLogger(logger log.Logger) Option {
	return func(l *Launcher) {
		l.logger = logger
	}
}
func WithKratosOptions(options ...kratos.Option) Option {
	return func(l *Launcher) {
		l.kratosOptions = append(l.kratosOptions, options...)
	}
}

func WithoutServiceDiscovery() Option {
	return func(l *Launcher) {
		l.notNeedServiceDiscovery = true
	}
}
