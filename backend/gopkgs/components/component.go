package components

import "github.com/go-kratos/kratos/v2/config"

// ConfigMap 是一个泛型映射，用于存储以字符串为键的任意类型配置项。
type ConfigMap[T any] map[string]T

// Component 定义了一个通用组件结构体，它能够根据配置项初始化组件，并提供错误处理机制。
// 它使用泛型来允许配置项可以是任何类型，增加了组件的灵活性和复用性。
type Component[T any] struct {
	err         error                                        // err 用于存储组件初始化过程中可能发生的错误。
	configValue map[string]config.Value                      // configValue 存储了组件的配置项，这些配置项可能还未被解析或应用。
	cfg         ConfigMap[T]                                 // cfg 是组件的当前配置，已经解析并准备好使用的配置项。
	initMethod  func(cfg ConfigMap[T]) (func() error, error) // initMethod 是一个函数，用于根据给定的配置初始化组件。它返回一个用于启动组件的函数和可能的初始化错误。
}

// Load 是一个泛型函数，用于加载配置并初始化组件。
// 它接受一个配置值的映射和一个初始化方法作为参数。
// 配置值的映射包含了所有需要加载的配置项，而初始化方法则定义了如何使用这些配置项来初始化组件。
// 初始化方法返回一个错误处理函数和一个可能的错误，这允许在初始化组件时执行一些延迟的错误处理逻辑。
func Load[T any](configValue map[string]config.Value, initMethod func(cfg ConfigMap[*T]) (func() error, error)) (t *T, components *Component[*T]) {
	// 创建一个组件实例，将配置值和初始化方法传递给它。
	c := &Component[*T]{
		configValue: configValue,
		initMethod:  initMethod,
	}
	// 初始化组件的配置映射。
	c.cfg = make(ConfigMap[*T])

	// 遍历配置值的映射，将每个配置项加载到组件的配置映射中。
	for key, value := range configValue {
		// 为每个配置项创建一个新的类型实例。
		t = new(T)
		// 使用Scan方法将配置值加载到类型实例中。
		// 如果加载失败，抛出一个panic。
		if err := value.Scan(t); err != nil {
			panic("scan componet config error: " + err.Error())
		}
		// 将加载的配置值存储在组件的配置映射中。
		c.cfg[key] = t
	}

	// 返回最后一个加载的类型实例和组件实例。
	return t, c
}

func (c *Component[T]) Start() error {
	if c.err != nil {
		return c.err
	}
	healthCheckMethod, err := c.initMethod(c.cfg)
	if err != nil {
		return err
	}
	return healthCheckMethod()
}
