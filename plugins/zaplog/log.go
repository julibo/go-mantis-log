package zaplog

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/julibo/go-mantis-log/conf"
	"github.com/julibo/go-mantis-log/fileout"
)

type Log struct {
	logger *zap.Logger
}

// var Log *zap.Logger // 全局日志

func parseLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "panic", "dpanic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	case "error":
		return zapcore.ErrorLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.DebugLevel

	}
}

// 创建日志
func New(opts ...conf.Option) *Log {
	o := &conf.Options{
		LogPath:     conf.LogPath,
		LogName:     conf.LogName,
		LogLevel:    conf.LogLevel,
		MaxSize:     conf.MaxSize,
		MaxAge:      conf.MaxAge,
		MaxBackups:  conf.MaxBackups,
		Stacktrace:  conf.Stacktrace,
		IsStdOut:    conf.IsStdOut,
		ProjectName: conf.ProjectName,
	}

	for _, opt := range opts {
		opt(o)
	}

	writers := []zapcore.WriteSyncer{fileout.NewRollingFile(o.LogPath, o.LogName, o.MaxSize, o.MaxAge, o.MaxBackups)}
	if o.IsStdOut == "yes" {
		writers = append(writers, os.Stdout)
	}

	logger := newZapLogger(parseLevel(o.LogLevel), parseLevel(o.Stacktrace), zapcore.NewMultiWriteSyncer(writers...))
	zap.RedirectStdLog(logger)
	logger = logger.With(zap.String("project", o.ProjectName))
	return &Log{logger: logger}
}

func newZapLogger(level, stacktrace zapcore.Level, output zapcore.WriteSyncer) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "@timestamp",
		NameKey:       "app",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		// EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeCaller: zapcore.ShortCallerEncoder, // 路径编码器
		EncodeName:   zapcore.FullNameEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-10-02 15:05:05.000"))
		},
	}

	// 设置日志级别
	dyn := zap.NewAtomicLevel()
	dyn.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg), // 编码器配置
		output,                         // 打印到文件和控制台
		dyn,                            // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "serviceName"))

	return zap.New(core, caller, development, zap.AddStacktrace(stacktrace), zap.AddCallerSkip(2), filed)
}
