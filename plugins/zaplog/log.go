package zaplog

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go-mantis-log/conf"
	"go-mantis-log/fileout"
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
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "@timestamp",
		NameKey:        "app",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		//EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		//	enc.AppendString(t.Format("2006-10-02 15:05:05.000"))
		//},
	}

	var encoder zapcore.Encoder
	// 设置日志级别
	dyn := zap.NewAtomicLevel()
	//encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	//encoder = zapcore.NewJSONEncoder(encCfg) // zapcore.NewConsoleEncoder(encCfg)
	dyn.SetLevel(level)
	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoder = zapcore.NewJSONEncoder(encCfg)

	return zap.New(zapcore.NewCore(encoder, output, dyn), zap.AddCaller(), zap.AddStacktrace(stacktrace), zap.AddCallerSkip(2))

}
