package defaultlogger

import (
	"os"

	"github.com/IsaacIsRebooting/Metok/backend/gopkgs/internal/defaultmiddleware"
	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LogPath       = "./logs/"
	LogFileName   = "log.log"
	LogMaxSize    = 5
	LogMaxBackups = 3
	LogMaxAge     = 30
	// LogCompress 指示是否在滚动时压缩日志文件备份。
	LogCompress = false
)

// getJsonLogEncoder 返回一个用于日志记录的 JSON 编码器。
// 该编码器根据特定的配置生成 JSON 格式的日志信息。
func getJsonLogEncoder() zapcore.Encoder {
	// 创建一个生产环境的编码器配置。
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置时间编码格式为 ISO8601，以确保时间信息在日志中的统一和可读性。
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置时间字段的键名为 "time"，以便在日志条目中标识时间信息。
	encoderConfig.TimeKey = "time"

	// 设置日志级别编码格式为大写形式，如 "INFO", "ERROR"，以统一日志级别格式。
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置持续时间编码格式为秒，简化持续时间的读取和理解。
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	// 设置调用者编码格式为短格式，仅包括文件名和行号，精简日志信息。
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 根据上述配置创建并返回一个 JSON 编码器。
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogSyncWriter 创建并返回一个同步日志写入器。
// 该函数不接受任何参数。
// 返回值是一个实现了WriteSyncer接口的对象，可以将日志写入到多个目的地。
func getLogSyncWriter() zapcore.WriteSyncer {
	// 初始化Lumberjack logger，用于滚动日志文件。
	// Lumberjack logger配置了日志文件路径、最大文件大小、备份文件数量、最大保存天数和是否压缩备份文件。
	lumberJackLogger := &lumberjack.Logger{
		Filename:   LogPath + LogFileName, // 日志文件路径和文件名
		MaxSize:    LogMaxSize,            // 最大文件大小（MB）
		MaxBackups: LogMaxBackups,         // 最多保留的备份文件数量
		MaxAge:     LogMaxAge,             // 最长保留的备份文件时间（天）
		Compress:   LogCompress,           // 是否压缩备份文件
	}
	// 返回一个多路写入同步器，它将日志写入到Lumberjack logger和标准输出。lumberJackLogger实现了WriteSyncer 接口
	// 这样做是为了同时支持将日志写入到文件和控制台。
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(lumberJackLogger),
		zapcore.AddSync(os.Stdout),
	)
}

// GetLogger 返回一个配置好的日志记录器。
// 该函数初始化日志记录器的核心配置，包括日志的编码方式、输出目的地和日志级别，
// 并通过Tracing和Kratos框架的适配，使其支持分布式追踪和框架特定的日志字段。
func GetLogger() log.Logger {
	// 获取同步写入日志的接口，该接口确保日志可以被有效地写入到指定的目的地。
	writeSyncer := getLogSyncWriter()
	// 获取JSON格式的日志编码器，用于将日志条目序列化为JSON格式。
	encoder := getJsonLogEncoder()

	// 创建一个新的核心配置，结合了JSON编码器、同步写入器和调试级别的日志记录。
	core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	// 初始化一个新的Zap日志记录器，使用上述核心配置。
	z := zap.New(core)
	// 初始化Tracing服务器，准备接收追踪信息。
	tracing.Server()
	// 使用Kratos框架的适配器包装Zap日志记录器，以适应Kratos的日志接口。
	zapLogger := kratoszap.NewLogger(z)
	// 在日志记录器中添加默认的和追踪相关的字段，包括时间戳、调用者信息和追踪ID等。
	logger := log.With(
		zapLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
		"x_trace_id", defaultmiddleware.GetTraceId(),
		"x_span_id", defaultmiddleware.GetSpanId(),
	)
	// 返回配置好的日志记录器。
	return logger
}
