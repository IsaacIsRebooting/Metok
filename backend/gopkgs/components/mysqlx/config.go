package mysqlx

import "fmt"

// Config 数据库配置结构体
// 该结构体用于定义数据库连接和操作的相关配置参数
type Config struct {
	Dialect           string `yaml:"dialect" json:"dialect"`                       // 数据库方言，如mysql、postgres等
	Host              string `yaml:"host" json:"host"`                             // 数据库主机地址
	Port              int    `yaml:"port" json:"port"`                             // 数据库端口号
	Dbname            string `yaml:"dbname" json:"dbname"`                         // 数据库名称
	User              string `yaml:"user" json:"user"`                             // 数据库用户名
	Password          string `yaml:"password" json:"password"`                     // 数据库密码
	Charset           string `yaml:"charset" json:"charset"`                       // 字符集
	ParseTime         bool   `yaml:"parse_time" json:"parse_time"`                 // 是否解析时间
	MaxIdle           int    `yaml:"max_idle" json:"max_idle"`                     // 最大空闲连接数
	MaxOpen           int    `yaml:"max_open" json:"max_open"`                     // 最大活跃连接数
	ConnMaxLifeTime   int    `yaml:"conn_max_life_time" json:"conn_max_life_time"` // 连接最大生命周期
	ConnMaxIdleTime   int    `yaml:"conn_max_idle_time" json:"conn_max_idle_time"` // 连接最大空闲时间
	Debug             bool   `yaml:"debug" json:"debug"`                           // 是否开启调试模式
	NoLog             bool   `yaml:"no_log" json:"no_log"`                         // 是否禁用日志记录
	InterpolateParams bool   `yaml:"interpolate_params" json:"interpolate_params"` // 是否插值SQL参数
	MultiStatements   bool   `yaml:"multi_statements" json:"multi_statements"`     // 是否支持多语句查询
	Timeout           int    `yaml:"timeout" json:"timeout"`                       // 连接超时时间（秒）
	ReadTimeout       int    `yaml:"read_timeout" json:"read_timeout"`             // 读取超时时间（秒）
	WriteTImeout      int    `yaml:"write_timeout" json:"write_timeout"`           // 写入超时时间（秒）
}

func (c *Config) setDefault() {
	if c.Dialect == "" {
		c.Dialect = "mysql"
	}

	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
	if !c.ParseTime {
		c.ParseTime = true
	}
	if c.MaxIdle == 0 {
		c.MaxIdle = 10
	}
	if c.MaxOpen == 0 {
		c.MaxOpen = 100
	}
}

func (c *Config) ToDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Dbname,
		c.Charset,
		c.ParseTime,
	)
	if c.InterpolateParams {
		dsn += "&interpolateParams=true"
	}
	if c.MultiStatements {
		dsn += "&multiStatements=true"
	}
	if c.Timeout > 0 {
		dsn += fmt.Sprintf("&timeout=%ds", c.Timeout)
	}
	if c.ReadTimeout > 0 {
		dsn += fmt.Sprintf("&readTimeout=%ds", c.ReadTimeout)
	}
	if c.WriteTImeout > 0 {
		dsn += fmt.Sprintf("&writeTimeout=%ds", c.WriteTImeout)
	}
	return dsn
}
