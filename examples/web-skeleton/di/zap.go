package di

import (
	"fmt"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xdi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func init() {
	obj := xdi.Object{
		Name: "zap",
		New: func() (i interface{}, e error) {
			filename := fmt.Sprintf("%s/../runtime/logs/mix.log", xcli.App().BasePath)
			fileRotate := &lumberjack.Logger{
				Filename:   filename,
				MaxBackups: 7,
			}
			atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
			core := zapcore.NewCore(
				zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
					TimeKey:       "T",
					LevelKey:      "L",
					NameKey:       "N",
					CallerKey:     "C",
					MessageKey:    "M",
					StacktraceKey: "S",
					LineEnding:    zapcore.DefaultLineEnding,
					EncodeLevel:   zapcore.CapitalLevelEncoder,
					EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
						enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
					},
					EncodeDuration: zapcore.StringDurationEncoder,
					EncodeCaller:   zapcore.ShortCallerEncoder,
				}),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileRotate)),
				atomicLevel,
			)
			logger := zap.New(core, zap.AddCaller())
			if xcli.App().Debug {
				atomicLevel.SetLevel(zap.DebugLevel)
			}
			return logger.Sugar(), nil
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}

func Zap() (logger *zap.SugaredLogger) {
	if err := xdi.Populate("zap", &logger); err != nil {
		panic(err)
	}
	return
}

type ZapOutput struct {
	Logger *zap.SugaredLogger
}

func (t *ZapOutput) Write(p []byte) (n int, err error) {
	t.Logger.Info(string(p))
	return len(p), nil
}
