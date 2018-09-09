// logRotator.go
// leoly
// 2017-07-23

package logRotator

/*
	日志文件拆分写
*/

import (
	"fmt"
	"log"
	"os"
	"time"
)

// StartLogRotator ...
func StartLogRotator(baseFileName string, maxSize int64, logger *log.Logger, rotatorMinute []int, timeRotate bool) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			fmt.Printf("StartLogRotator panic:[%v]\n", panicErr)
		}
	}()
	LogBaseName = baseFileName
	LogMaxSize = maxSize
	Logger = logger
	if nil != rotatorMinute && len(rotatorMinute) > 0 {
		LogRotatorMinute = rotatorMinute
	} else {
		NeedCheckTime = timeRotate
	}

	initLog()

	go func() {
		for {
			fileChecker()
			time.Sleep(time.Second * time.Duration(LogCheckInterval))
		}
	}()
}

func initLog() error {
	outputFile, err := openFile(LogBaseName)
	if nil == outputFile {
		return err
	}
	setLoggerOutput(outputFile)
	LoggerFileObj = outputFile
	return nil
}

func setLoggerOutput(outputFile *os.File) {
	if nil == Logger {
		log.SetOutput(outputFile)
	} else {
		Logger.SetOutput(outputFile)
	}
}

// Logger ...
var Logger *log.Logger

// LoggerFileObj ...
var LoggerFileObj *os.File

// LogBaseName ...
var LogBaseName string // 文件名

// LogMaxSize ...
var LogMaxSize int64 = 512 * 1024 * 1024 // 文件大小, 默认 500MB

// LogCheckInterval ...
var LogCheckInterval int64 = 5 // 检查周期 默认 1s

// LogRotatorMinute ...
var LogRotatorMinute = []int{0, 30} // 滚动分钟, 默认整点,半点

// NeedCheckTime ...
var NeedCheckTime = true // 是否需要检测切换时间

func openFile(baseFileName string) (*os.File, error) {
	monitorLogFile := fmt.Sprintf("%v_%v", baseFileName, time.Now().Format("20060102150405.000000"))

	outputFile, err := os.OpenFile(monitorLogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if nil != err {
		fmt.Printf("open log file [%v] failed:[%v]\n", monitorLogFile, err)
		return nil, err
	}
	return outputFile, nil
}

func fileChecker() {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			fmt.Printf("fileChecker() panic:[%v]\n", panicErr)
		}
	}()

	stat, err := LoggerFileObj.Stat()
	fileSize := LogMaxSize // default fileSize > LogMaxSize, if stat failed, roll file
	if nil != err {
		fmt.Printf("get file stat failed:%v\n", err)
	} else {
		fileSize = stat.Size()
	}
	// fileName := stat.Name()
	// fmt.Printf("fileName:[%v], size:[%v]\n", fileName, fileSize)

	if checkSize(fileSize) || checkTime(time.Now()) {
		outputFile, err := openFile(LogBaseName)
		if err == nil {
			setLoggerOutput(outputFile)
			LoggerFileObj.Close()
			LoggerFileObj = outputFile
		}
	}
}

// check
//	true: need rotate file
//	false: do not need rotate file
func checkSize(fileSize int64) bool {
	if fileSize >= LogMaxSize {
		return true
	}
	return false
}

func checkTime(timeNow time.Time) bool {
	if !NeedCheckTime {
		return false
	}
	if timeNow.Second() == 0 {
		for _, minute := range LogRotatorMinute {
			if timeNow.Minute() == minute {
				return true
			}
		}
	}
	return false
}
