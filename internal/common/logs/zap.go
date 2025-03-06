package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
	"xmicro/internal/common/config/center"
)

var (
	logger     *zap.SugaredLogger
	initLogger sync.Once
)

func Init() {
	initLogger.Do(func() {
		writeSyncer := getLogWriter(config.Conf.ZapLog.File)
		encoder := getEncoder()
		core := zapcore.NewCore(encoder, writeSyncer, config.Conf.ZapLog.Level)

		opts := []zap.Option{zap.AddCaller()}
		if config.LocalConf.AppName != "" {
			opts = append(opts, zap.Fields(zap.String("app", config.LocalConf.AppName)))
		}

		log := zap.New(core, opts...)
		logger = log.Sugar()
	})
	if logger == nil {
		panic("failed to initialize logger")
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logFile string) zapcore.WriteSyncer {
	if logFile == "" {
		logFile = "./test.logs" // 默认日志文件路径
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}

	ws := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(ws)
}

func logMessage(level zapcore.Level, format string, values ...interface{}) {
	msg := fmt.Sprintf(format, values...)
	switch level {
	case zapcore.DebugLevel:
		logger.Debug(msg)
	case zapcore.InfoLevel:
		logger.Info(msg)
	case zapcore.WarnLevel:
		logger.Warn(msg)
	case zapcore.ErrorLevel:
		logger.Error(msg)
	case zapcore.FatalLevel:
		logger.Fatal(msg)
	default:
		logger.Info(msg)
	}
}

func Debug(format string, values ...interface{}) {
	logMessage(zapcore.DebugLevel, format, values...)
}

func Info(format string, values ...interface{}) {
	logMessage(zapcore.InfoLevel, format, values...)
}

func Warn(format string, values ...interface{}) {
	logMessage(zapcore.WarnLevel, format, values...)
}

func Error(format string, values ...interface{}) {
	logMessage(zapcore.ErrorLevel, format, values...)
}

func Fatal(format string, values ...interface{}) {
	logMessage(zapcore.FatalLevel, format, values...)
}
