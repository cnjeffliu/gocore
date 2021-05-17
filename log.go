/*
* 日志组件封装
* jeff.liu <zhifeng172@163.com> 20190126
 */
package util

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLogger *zap.Logger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func toStr(template string, fmtArgs []interface{}) string {
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}
	return msg
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%6s", l.CapitalString()))
}

// specify codeline's filename and row
func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fullPath := caller.FullPath()
	// idx := strings.LastIndexByte(fullPath, ':')
	// num := uint64(0)
	// if idx != -1 {
	// 	num, _ = strconv.ParseUint(fullPath[idx+1:], 10, 32)
	// 	fullPath = fullPath[:idx]
	// }

	// if len(fullPath) > 20 {
	// 	fullPath = fullPath[len(fullPath)-20:]
	// }

	// enc.AppendString(fmt.Sprintf("%20s:%d", fullPath, num))

	if len(fullPath) > 20 {
		enc.AppendString(fmt.Sprintf("%20.20s", fullPath[len(fullPath)-20:]))
	} else {
		enc.AppendString(fmt.Sprintf("%20.20s", fullPath))
	}
}

var ALevel zap.AtomicLevel

func InitLog(filename string) {
	ALevel = zap.NewAtomicLevel()
	ALevel.SetLevel(getLoggerLevel("debug"))

	hook := lumberjack.Logger{
		Filename:   filename,
		MaxSize:    30, // megabytes
		MaxBackups: 10,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
	}

	fileWriter := zapcore.AddSync(&hook)

	//consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	//core := zapcore.NewTee(
	//	// 打印在文件中
	//	zapcore.NewCore(consoleEncoder, fileWriter, highPriority),
	//)
	//zapLogger = zap.New(core)

	// config := zap.NewProductionEncoderConfig()
	config := zap.NewDevelopmentEncoderConfig()
	config.ConsoleSeparator = " "

	config.EncodeLevel = levelEncoder
	config.EncodeTime = timeEncoder
	config.EncodeCaller = callerEncoder

	//encoder := zapcore.NewJSONEncoder(config)
	encoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewCore(encoder, fileWriter, ALevel)
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func SetLevel(lvl string) {
	ALevel.SetLevel(getLoggerLevel(lvl))
}

func Debug(args ...interface{}) {
	zapLogger.Debug(fmt.Sprint(args...))
}

func Debugf(fmt string, args ...interface{}) {
	zapLogger.Debug(toStr(fmt, args))
}

func Info(args ...interface{}) {
	zapLogger.Info(fmt.Sprint(args...))
}

func Infof(fmt string, args ...interface{}) {
	zapLogger.Info(toStr(fmt, args))
}

func Warn(args ...interface{}) {
	zapLogger.Warn(fmt.Sprint(args...))
}

func Warnf(fmt string, args ...interface{}) {
	zapLogger.Warn(toStr(fmt, args))
}

func Error(args ...interface{}) {
	zapLogger.Error(fmt.Sprint(args...))
}
func Errorf(fmt string, args ...interface{}) {
	zapLogger.Error(toStr(fmt, args))
}

func Fatal(args ...interface{}) {
	zapLogger.Fatal(fmt.Sprint(args...))
}
func Fatalf(fmt string, args ...interface{}) {
	zapLogger.Fatal(toStr(fmt, args))
}
