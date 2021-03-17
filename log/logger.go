package log

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/goava/di"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/urionz/goofy/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
	FatalLevel = "fatal"
)

type Logger struct {
	di.Tags `name:"logger"`

	*zap.Logger
	conf contracts.Config
	app  contracts.Application
}

var log *Logger

func Debug(args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Errorf(format, args...)
}

func Panic(args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalf(format, args...)
}

func NewLogger(app contracts.Application, conf contracts.Config) *Logger {
	log = new(Logger)
	log.conf = conf
	log.app = app
	log.Logger = log.newZapLogger(log.parseLogLevel(conf.String("logger.level", "debug")))
	return log
}

func (logger *Logger) newZapLogger(level zapcore.Level) *zap.Logger {
	conf := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(logger.conf.String("logger.encode_time", "2006-01-02 15:04:05")))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if logger.conf.Bool("logger.color", true) {
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	encoder := zapcore.NewConsoleEncoder(conf)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= level
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= level
	})

	infoLevelWriter := logger.getLevelWriter("./info")
	warnLevelWriter := logger.getLevelWriter("./warn")

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoLevelWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnLevelWriter), warnLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(conf), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), level),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func (logger *Logger) getLevelWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		path.Join(logger.app.Storage(), filename+"-%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join(logger.app.Storage(), filename+".log")),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func (*Logger) parseLogLevel(level string) zapcore.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case PanicLevel:
		return zap.PanicLevel
	case FatalLevel:
		return zap.FatalLevel
	}
	return zap.InfoLevel
}
