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

type Logger struct {
	di.Tags `name:"logger"`

	*zap.Logger
	atomic zap.AtomicLevel
	conf   contracts.Config
	app    contracts.Application
	output string
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

func Sugar() *zap.SugaredLogger {
	return log.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func GetRotateWriter(filename string) io.Writer {
	return log.getRotateWriter(filename)
}

func NewLogger(app contracts.Application, conf contracts.Config) *Logger {
	log = new(Logger)
	log.conf = conf
	log.app = app
	output := path.Join(app.Storage(), conf.String("logger.output_path", "logs"))
	absOutput := conf.String("logger.output_path_abs")
	if absOutput != "" {
		output = absOutput
	}
	log.output = output
	log.Logger = log.newZapLogger()
	return log
}

func (logger *Logger) newZapLogger() *zap.Logger {
	var encoder zapcore.Encoder
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

	if logger.conf.String("app.env", "production") == "production" {
		conf.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(conf)
	} else {
		encoder = zapcore.NewConsoleEncoder(conf)
	}

	logger.atomic = zap.NewAtomicLevel()

	logger.DynamicConf(logger.app, logger.conf)

	coreTee := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(logger.getRotateWriter("logger")),
			zapcore.AddSync(os.Stdout),
		), logger.atomic),
	}

	if logger.conf.Bool("logger.multi_level_output", true) {
		infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.WarnLevel && lvl >= logger.atomic.Level()
		})

		infoLevelWriter := logger.getRotateWriter("info")

		warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.WarnLevel && lvl >= logger.atomic.Level()
		})
		warnLevelWriter := logger.getRotateWriter("warn")
		coreTee = append(coreTee, zapcore.NewCore(encoder, zapcore.AddSync(infoLevelWriter), infoLevel))
		coreTee = append(coreTee, zapcore.NewCore(encoder, zapcore.AddSync(warnLevelWriter), warnLevel))
	}

	core := zapcore.NewTee(coreTee...)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func (logger *Logger) DynamicConf(_ contracts.Application, conf contracts.Config) error {
	logger.SetLevel(conf.String("logger.level", "debug"))
	return nil
}

func (logger *Logger) SetLevel(level string) {
	logger.atomic.SetLevel(logger.parseLogLevel(level))
}

func (logger *Logger) getRotateWriter(filename string) io.Writer {
	maxAge, _ := time.ParseDuration(logger.conf.String("logger.rotate.max_age", "240h"))
	period, _ := time.ParseDuration(logger.conf.String("logger.rotate.period", "24h"))
	filename = path.Join(logger.output, filename)
	hook, err := rotatelogs.New(
		filename+"-%Y-%m-%d.log",
		rotatelogs.WithLinkName(filename+".log"),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(period),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func (*Logger) parseLogLevel(level string) zapcore.Level {
	switch level {
	case contracts.DebugLevel:
		return zap.DebugLevel
	case contracts.ErrorLevel:
		return zap.ErrorLevel
	case contracts.InfoLevel:
		return zap.InfoLevel
	case contracts.WarnLevel:
		return zap.WarnLevel
	case contracts.PanicLevel:
		return zap.PanicLevel
	case contracts.FatalLevel:
		return zap.FatalLevel
	}
	return zap.InfoLevel
}
