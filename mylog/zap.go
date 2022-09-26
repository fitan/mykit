package mylog

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
)

func New(name string, dir string, level zap.AtomicLevel) *zap.SugaredLogger {
	l := zap.New(DefaultZapCore(name, dir, level)).Sugar()
	return l
}

func DefaultZapCore(fileName string, dir string, openLevel zap.AtomicLevel) zapcore.Core {
	errEnable := zap.LevelEnablerFunc(
		func(level zapcore.Level) bool {
			return level >= zap.ErrorLevel && zap.ErrorLevel >= openLevel.Level()
		})
	infoEnable := zap.LevelEnablerFunc(
		func(level zapcore.Level) bool {
			return level == zap.InfoLevel && level >= openLevel.Level()
		})
	warnEnable := zap.LevelEnablerFunc(
		func(level zapcore.Level) bool {
			return level == zap.WarnLevel && level >= openLevel.Level()
		})
	debugEnable := zap.LevelEnablerFunc(
		func(level zapcore.Level) bool {
			return level == zap.DebugLevel && level >= openLevel.Level()
		})

	infoLogWriter := getLogWriter(path.Join(dir, fileName+"_info.log"))
	errLogWriter := getLogWriter(path.Join(dir, fileName+"_err.log"))
	warnWriter := getLogWriter(path.Join(dir, fileName+"_warn.log"))
	debugWriter := getLogWriter(path.Join(dir, fileName+"_debug.log"))

	stdout := zapcore.AddSync(colorable.NewColorableStdout())
	infoCore := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(infoLogWriter, stdout), infoEnable)
	errCore := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(errLogWriter, stdout), errEnable)
	warnCore := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(warnWriter, stdout), warnEnable)
	debugCore := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(debugWriter, stdout), debugEnable)

	return zapcore.NewTee(infoCore, errCore, warnCore, debugCore)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
	//return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
