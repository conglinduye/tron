// package level log, inner use only

package log

import (
	"fmt"
	golog "log"
	"strings"
)

//Level 日志的级别
type Level int

//日志的级别定义
const (
	ALL Level = iota
	MORE
	SQL
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

//默认日志级别
const (
	defaultLogLevel = INFO
)

//isLogLevelValid 判断loglevel是否是一个有效值
func isLogLevelValid(level Level) bool {
	if ALL <= level && level <= OFF {
		return true
	}
	return false
}

//Str2Level 字符转LogLevel
func Str2Level(level string) (ret Level) {
	ret = INFO

	level = strings.ToLower(level)
	switch level {
	case "all":
		ret = ALL
	case "more":
		ret = MORE
	case "sql":
		ret = SQL
	case "debug":
		ret = DEBUG
	case "info":
		ret = INFO
	case "warn":
		ret = WARN
	case "error":
		ret = ERROR
	case "fatal":
		ret = FATAL
	case "off":
		ret = OFF
	default:
		ret = INFO
	}
	return ret
}

//isLogShouldRecord 判断日志是否应该被记录
func isLogShouldRecord(level Level) bool {
	//fmt.Printf("currentlogLevel:[%v] level:[%v]\n", currentlogLevel, level)
	return level >= currentlogLevel
}

//当前的日志等级
var currentlogLevel = Level(defaultLogLevel)

//ChangeLogLevel 更改当前的日志等级
func ChangeLogLevel(level Level) bool {
	if isLogLevelValid(level) {
		var oldLevel = currentlogLevel
		currentlogLevel = level
		Infof("change log level from [%v] to [%v] ", oldLevel, currentlogLevel)

		return true
	}
	return false
}

// Output 底层的输出接口，可以在此之上进行封装，直接调用时 stacklevel == 3 可以正常输出 stack position
//	每封装一层需要自己以一
func Output(stackLevel int, msg string) {
	golog.Output(stackLevel+1, msg)
}

func Print(args ...interface{}) {
	golog.Print(args...)
}

func Printf(format string, args ...interface{}) {
	golog.Println(args...)
}

func Println(args ...interface{}) {
	golog.Println(args...)
}

func Fatal(args ...interface{}) {
	golog.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	golog.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	golog.Fatalln(args...)
}

func Panic(args ...interface{}) {
	golog.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	golog.Panicf(format, args...)
}

func Panicln(args ...interface{}) {
	golog.Panicln(args...)
}

func More(args ...interface{}) {
	outputLogMessage(MORE, " [MORE] "+fmt.Sprint(args...))
}

func Moref(format string, args ...interface{}) {
	outputLogMessage(MORE, " [MORE] "+fmt.Sprintf(format, args...))
}

func Moreln(args ...interface{}) {
	outputLogMessage(MORE, " [MORE] "+fmt.Sprintln(args...))
}

func Sql(args ...interface{}) {
	outputLogMessage(SQL, " [SQL] "+fmt.Sprint(args...))
}

func Sqlf(format string, args ...interface{}) {
	outputLogMessage(SQL, " [SQL] "+fmt.Sprintf(format, args...))
}

func Sqlln(args ...interface{}) {
	outputLogMessage(SQL, " [SQL] "+fmt.Sprintln(args...))
}

func Debug(args ...interface{}) {
	outputLogMessage(DEBUG, " [DEBUG] "+fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	outputLogMessage(DEBUG, " [DEBUG] "+fmt.Sprintf(format, args...))
}

func Debugln(args ...interface{}) {
	outputLogMessage(DEBUG, " [DEBUG] "+fmt.Sprintln(args...))
}

func Info(args ...interface{}) {
	outputLogMessage(INFO, " [INFO] "+fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	outputLogMessage(INFO, " [INFO] "+fmt.Sprintf(format, args...))
}

func Infoln(args ...interface{}) {
	outputLogMessage(INFO, " [INFO] "+fmt.Sprintln(args...))
}

func Warn(args ...interface{}) {
	outputLogMessage(WARN, " [WARN] "+fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	outputLogMessage(WARN, " [WARN] "+fmt.Sprintf(format, args...))
}

func Warnln(args ...interface{}) {
	outputLogMessage(WARN, " [WARN] "+fmt.Sprintln(args...))
}

func Error(args ...interface{}) {
	outputLogMessage(ERROR, " [ERROR] "+fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	outputLogMessage(ERROR, " [ERROR] "+fmt.Sprintf(format, args...))
}

func Errorln(args ...interface{}) {
	outputLogMessage(ERROR, " [ERROR] "+fmt.Sprintln(args...))
}

//outputLogMessage 真实的日志输出
func outputLogMessage(level Level, message string) {
	if isLogShouldRecord(level) {
		err := golog.Output(3, message)
		if err != nil {
			fmt.Println("*** log message error. ***", err)
		}
	}
}
