package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"server-template/internal/conf"
	"server-template/pkg/env"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(logConfig *conf.Log, version, env string) {
	var err error
	var hostname string

	if hostname, err = os.Hostname(); err != nil {
		hostname = "localhost"
	}

	logger := New(
		os.Stderr, logConfig, env, WithCaller(true), AddCallerSkip(4),
		Fields(
			Field{
				Key:    "app",
				Type:   zapcore.StringType,
				String: logConfig.AppName,
			},
			Field{
				Key:    "host",
				Type:   zapcore.StringType,
				String: hostname,
			},
			Field{
				Key:    "version",
				Type:   zapcore.StringType,
				String: version,
			},
		),
	)

	ResetDefault(logger)
}

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	l      *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level  zap.AtomicLevel
	msgKey string
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.Debug(msg, data...)
	case log.LevelInfo:
		l.Info(msg, data...)
	case log.LevelWarn:
		l.Warn(msg, data...)
	case log.LevelError:
		l.Error(msg, data...)
	case log.LevelFatal:
		l.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Debugf(msg string, args ...any) {
	l.l.Sugar().Debugf(msg, args)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Infof(msg string, args ...any) {
	l.l.Sugar().Infof(msg, args...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Warnf(msg string, args ...any) {
	l.l.Sugar().Warnf(msg, args...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) Errorf(msg string, args ...any) {
	l.l.Sugar().Errorf(msg, args...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) DPanicf(msg string, args ...any) {
	l.l.Sugar().DPanicf(msg, args...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Panicf(msg string, args ...any) {
	l.l.Sugar().Panicf(msg, args...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(msg string, args ...any) {
	l.l.Sugar().Fatalf(msg, args...)
}

func (l *Logger) With(fields ...Field) *Logger {
	logger := l.l.With(fields...)
	return &Logger{l: logger, level: l.level, msgKey: l.msgKey}
}

var (
	std = New(os.Stderr, &conf.Log{}, env.ModeProd, WithCaller(true), AddCallerSkip(1))

	Info    = std.Info
	Infof   = std.Infof
	Warn    = std.Warn
	Warnf   = std.Warnf
	Error   = std.Error
	Errorf  = std.Errorf
	DPanic  = std.DPanic
	DPanicf = std.DPanicf
	Panic   = std.Panic
	Panicf  = std.Panicf
	Fatal   = std.Fatal
	Fatalf  = std.Fatalf
	Debug   = std.Debug
	Debugf  = std.Debugf
	With    = std.With

	mutex = sync.Mutex{}
)

// safe for concurrent use, replace default std
func ResetDefault(l *Logger) {
	mutex.Lock()
	defer mutex.Unlock()

	std = l
	Info = std.Info
	Infof = std.Infof
	Warn = std.Warn
	Warnf = std.Warnf
	Error = std.Error
	Errorf = std.Errorf
	DPanic = std.DPanic
	DPanicf = std.DPanicf
	Panic = std.Panic
	Panicf = std.Panicf
	Fatal = std.Fatal
	Fatalf = std.Fatalf
	Debug = std.Debug
	Debugf = std.Debugf
	With = std.With
}

func Default() *Logger { return std }

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddCallerSkip = zap.AddCallerSkip
	AddStacktrace = zap.AddStacktrace
	Fields        = zap.Fields
)

// New create a new logger
func New(w io.Writer, logConfig *conf.Log, envArg string, opts ...Option) *Logger {
	if w == nil {
		panic("the writer is nil")
	}

	encConfig := zap.NewProductionEncoderConfig()
	// set time format
	encConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encConfig.MessageKey = "message"

	atomicLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	var enc zapcore.Encoder
	switch envArg {
	case env.ModeDev, env.ModeTest:
		enc = zapcore.NewConsoleEncoder(encConfig)
		atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case env.ModeEmpty, env.ModeProd:
		enc = zapcore.NewJSONEncoder(encConfig)
	}

	var cores []zapcore.Core
	if logConfig.IsWriteFile {
		fileWriterSync := getFileLogWriter(
			logConfig.LogFile.Name, false,
			int(logConfig.LogFile.MaxSize), int(logConfig.LogFile.MaxAge),
		)
		core := zapcore.NewCore(
			enc,
			zapcore.AddSync(fileWriterSync),
			atomicLevel,
		)
		cores = append(cores, core)
	}

	core := zapcore.NewCore(
		enc,
		zapcore.AddSync(w),
		atomicLevel,
	)

	cores = append(cores, core)

	core = zapcore.NewTee(cores...)
	return &Logger{
		l:      zap.New(core, opts...),
		level:  atomicLevel,
		msgKey: log.DefaultMessageKey,
	}
}

// SetLevel alters the logging level on runtime
// it is concurrent-safe
func (l *Logger) SetLevel(level Level) {
	l.level.SetLevel(zapcore.Level(level))
}

func SetLevel(level Level) {
	std.level.SetLevel(zapcore.Level(level))
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}

	return nil
}

func getFileLogWriter(
	filePath string, isCompress bool,
	fileMaxSize, backUpFileMaxAge int,
) (writeSyncer zapcore.WriteSyncer) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  filePath,
		MaxSize:   fileMaxSize,
		MaxAge:    backUpFileMaxAge,
		Compress:  false,
		LocalTime: true,
	}

	return zapcore.AddSync(lumberJackLogger)
}
