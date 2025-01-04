package launcher

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/components/consulx"
	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/internal/defaultlogger"
	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/internal/shutdown"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Launcher is a structure that encapsulates the necessary components and configurations for starting an application.
// It includes configurations, logger, server setup, and various lifecycle handlers.
type Launcher struct {
	// app is an instance of the Kratos application framework, used to manage the application's lifecycle.
	app *kratos.App

	// notNeedServiceDiscovery indicates whether service discovery is needed. If true, service discovery is not required for this application.
	notNeedServiceDiscovery bool

	// configOptions contains the configuration options, used to initialize the application configuration.
	configOptions []config.Option

	// configWatchMap is a map that holds configuration observers, used to monitor changes in configuration.
	configWatchMap map[string]config.Observer

	// config represents the application's configuration, obtained through configuration initialization.
	config config.Config

	// configValue is used to store the parsed configuration values, the specific type depends on the application's configuration structure.
	configValue interface{}

	// logger is the logger instance, used to output logs during application runtime.
	logger log.Logger

	// grpcServer is a function that returns a gRPC server instance based on the configuration value.
	grpcServer func(configValue interface{}) *grpc.Server

	// ginServer is a function that returns a Gin HTTP server instance based on the configuration value.
	ginServer func(configValue interface{}) *http.Server

	// kratosOptions contains the options for initializing the Kratos application framework.
	kratosOptions []kratos.Option

	// componentsLauncher is used to manage the startup and shutdown of other component services.
	componentsLauncher *ComponentsLauncher

	// beforeConfigInitHandlers contains the handlers to be executed before configuration initialization.
	beforeConfigInitHandlers []func()

	// afterConfigInitHandlers contains the handlers to be executed after configuration initialization.
	afterConfigInitHandlers []func()

	// beforeServerStartHandlers contains the handlers to be executed before the server starts.
	beforeServerStartHandlers []func()

	// afterServerStartHandlers contains the handlers to be executed after the server starts.
	afterServerStartHandlers []func()

	// shutdownHandlers contains the handlers to be executed during application shutdown.
	shutdownHandlers []func()
}

// NewLauncher 创建并返回一个新的Launcher实例。
// 它接受一系列的Option作为参数，这些Option用于配置Launcher的属性。
// 这种设计模式允许在创建Launcher时灵活地设置各种选项，而不需要定义多个重载函数。
func NewLauncher(options ...Option) *Launcher {
	// 初始化Launcher实例，并创建一个用于观察配置变化的映射。
	l := &Launcher{
		configWatchMap: make(map[string]config.Observer),
	}

	// 遍历所有提供的Option函数，并对Launcher实例进行配置。
	for _, option := range options {
		option(l)
	}
	// 返回配置完毕的Launcher实例。
	return l
}

// Run 启动应用程序。
// 该方法主要负责应用程序的初始化、组件启动、服务器启动前后的处理以及关闭流程。
func (l *Launcher) Run() {
	// 获取当前工作目录并记录日志
	dir, _ := os.Getwd()
	log.Context(context.Background()).Info("current dir:", dir)

	// 初始化配置
	l.runInitConfig()

	// 执行服务器启动前的处理函数
	l.runHandlers(l.beforeServerStartHandlers, "start to run handlers before server start")

	// 启动组件
	l.componentsLauncher.Launch()

	// 创建新的Kratos应用实例
	l.newKratosApp()

	// 等待应用运行
	<-l.run()

	// 执行服务器启动后的处理函数,启动mysql
	l.runHandlers(l.afterServerStartHandlers, "start to run handlers after server start")

	// 等待关闭信号
	<-shutdown.FiredCh()

	// 在关闭前等待最多10秒以优雅关闭
	shutdown.Wait(10 * time.Second)

	// 执行服务器关闭时的处理函数
	l.runHandlers(l.shutdownHandlers, "start to run handlers during server shutdown")
}

// runInitConfig 初始化配置。
// 该方法首先执行预配置初始化处理程序，然后创建并加载配置对象。
// 加载成功后，它会使用新配置更新Launcher的配置属性，并为配置项添加观察者。
// 最后，执行后配置初始化处理程序。
func (l *Launcher) runInitConfig() {
	// 执行预配置初始化处理程序
	l.runHandlers(l.beforeConfigInitHandlers, "start to run handlers before config init")

	// 创建新的配置对象
	cfg := config.New(l.configOptions...)

	// 加载配置，如果失败则抛出错误
	if err := cfg.Load(); err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	// 更新Launcher的配置属性
	l.config = cfg

	// 为配置项添加观察者
	for key, observer := range l.configWatchMap {
		if err := cfg.Watch(key, observer); err != nil {
			panic(fmt.Errorf("failed to watch config: %v", err))
		}
	}

	// 将配置值扫描到Launcher的配置值属性中
	if err := cfg.Scan(l.configValue); err != nil {
		panic(fmt.Errorf("failed to scan config value: %v", err))
	}

	// 创建新的组件启动器，并更新Launcher的组件启动器属性
	l.componentsLauncher = NewComponentsLauncher(cfg)
	// 执行后配置初始化处理程序
	l.runHandlers(l.afterConfigInitHandlers, "start to run handlers after config init")
}

func (l *Launcher) initTracer(cfg *App) {
	if cfg.TraceEndpoint != "" {
		if err := initTracer(cfg.Name, cfg.TraceEndpoint); err != nil {
			panic(err)
		}
	}
}

// runHandlers 执行一组给定的处理器函数，并在执行前记录一条日志消息。
// 此函数的目的是在特定上下文中执行一系列操作，每个操作由一个处理器函数表示。
// 参数:
//
//	handlers - 一个包含多个处理器函数的切片，这些函数将被依次执行。
//	msg - 在执行任何处理器函数之前要记录的日志消息，用于提供执行上下文。
func (l *Launcher) runHandlers(handlers []func(), msg string) {
	// 检查是否有处理器函数需要执行，如果没有，则不执行任何操作。
	if len(handlers) > 0 {
		// 在执行处理器函数之前记录日志消息，提供执行这些处理器的上下文信息。
		log.Context(context.Background()).Info(msg)
	}

	// 遍历处理器函数切片，执行每个处理器函数。
	// 这里没有并行执行处理器函数，因为它们可能依赖于特定的执行顺序或共享资源。
	for _, handler := range handlers {
		handler()
	}
}

// newKratosApp 创建一个新的Kratos应用实例。
// 该方法主要用于初始化Kratos应用的配置，特别是日志选项。
func (l *Launcher) newKratosApp() {
	// 初始化一个空的Kratos选项切片。
	options := make([]kratos.Option, 0)

	// 如果Launcher实例已经设置了自定义日志记录器，
	// 则将其添加到Kratos应用的选项中。
	if l.logger != nil {
		options = append(options, kratos.Logger(l.logger))
	} else {
		// 如果没有设置自定义日志记录器，则使用默认的日志记录器。
		// 这里通过defaultlogger包获取默认的日志记录器实例并添加到选项中。
		options = append(options, kratos.Logger(defaultlogger.GetLogger()))
	}
	if l.grpcServer != nil {
		options = append(options, kratos.Server(l.grpcServer(l.configValue)))
	}
	if l.ginServer != nil {
		options = append(options, kratos.Server(l.ginServer(l.configValue)))
	}
	if len(l.kratosOptions) > 0 {
		options = append(options, l.kratosOptions...)
	}
	if !l.notNeedServiceDiscovery {
		consulClient := consulx.GetClient(context.Background())
		consulReg := consul.New(consulClient)
		options = append(options, kratos.Registrar(consulReg))
	}
	value := l.config.Value("app")
	appConfig := &App{}
	if err := value.Scan(appConfig); err != nil {
		panic(fmt.Errorf("failed to scan app config, error: %v", err))
	}
	options = append(options, kratos.Name(appConfig.Name), kratos.Version(appConfig.Version))

	l.app = kratos.New(options...)

	l.initTracer(appConfig)
}

// runKratosApp 启动一个Kratos应用，并在应用准备好后通过通道发送信号。
// 该函数首先检查应用实例是否已经初始化，如果没有则抛出panic。
// 返回一个通道，当应用成功启动并准备好时，该通道将被关闭。
func (l *Launcher) runKratosApp() <-chan struct{} {
	// 检查应用实例是否已经初始化
	if l.app == nil {
		panic("app not initialized")
	}

	// 创建一个通道，用于通知应用已准备好
	readyChan := make(chan struct{})

	// 初始化一个等待组，用于同步应用启动过程
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 启动应用实例
	go func() {
		wg.Done()

		// 如果应用启动失败，则记录错误并抛出panic
		if err := l.app.Run(); err != nil {
			log.Context(context.Background()).Fatal("failed to run app")
			panic(err)
		}
	}()

	// 等待应用实例启动，并关闭准备就绪通道
	go func() {
		wg.Wait()
		close(readyChan)
	}()

	// 返回准备就绪通道
	return readyChan
}

func (l *Launcher) run() <-chan struct{} {
	ch := l.runKratosApp()
	appReadyCh := make(chan struct{})

	go func() {
		<-ch
		close(appReadyCh)
	}()
	return appReadyCh
}
